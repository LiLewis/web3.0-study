package task4

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 配置 & 全局变量
var (
	// JWT 密钥（请在生产环境中从环境变量读取）
	jwtSecret = []byte(getEnv("JWT_SECRET", "dev_secret_change_me"))

	// 数据库句柄
	db *gorm.DB
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// User Post Comment 模型定义（GORM）
type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Username string `gorm:"uniqueIndex;size:64;not null"` // 登录名
	Password string `gorm:"not null"`                     // bcrypt 哈希
	Email    string `gorm:"uniqueIndex;size:128;not null"`

	// 统计字段（演示 Hook 用）
	PostCount int `gorm:"default:0"`
}

type Post struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title   string `gorm:"size:200;not null"`
	Content string `gorm:"type:text;not null"`

	UserID uint
	User   User

	// 评论关联 + 状态（演示 Hook 用）
	Comments      []Comment
	CommentStatus string `gorm:"size:20;default:'无评论'"`
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Content string `gorm:"type:text;not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

// AfterCreate --- Hooks ---
// 文章创建后，给作者的 PostCount +1
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).Error
}

// AfterCreate 评论创建后，将文章 CommentStatus 标为 "有评论"
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&Post{}).Where("id = ?", c.PostID).
		Update("comment_status", "有评论").Error
}

// AfterDelete 评论删除后，如该文章已无评论，则置为 "无评论"
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)
	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error
	}
	return nil
}

// RegisterReq LoginReq PostCreateReq PostUpdateReq CommentCreateReq 请求/响应 DTO
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostCreateReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PostUpdateReq struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type CommentCreateReq struct {
	Content string `json:"content" binding:"required"`
}

// 工具函数
func hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func jsonError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

// JWTClaims JWT 鉴权中间件
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func issueToken(u User) (string, error) {
	claims := JWTClaims{
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			jsonError(c, http.StatusUnauthorized, "missing token")
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			jsonError(c, http.StatusUnauthorized, "invalid token")
			return
		}
		claims := token.Claims.(*JWTClaims)
		c.Set("uid", claims.UserID)
		c.Set("uname", claims.Username)
		c.Next()
	}
}

// 只能作者本人操作文章
func authorOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uidVal, ok := c.Get("uid")
		if !ok {
			jsonError(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		var p Post
		id := c.Param("id")
		if err := db.First(&p, id).Error; err != nil {
			jsonError(c, http.StatusNotFound, "post not found")
			return
		}
		if p.UserID != uidVal.(uint) {
			jsonError(c, http.StatusForbidden, "forbidden: not the author")
			return
		}
		c.Next()
	}
}

// 处理器：认证
func registerHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	pw, err := hashPassword(req.Password)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "hash password failed")
		return
	}
	u := User{Username: req.Username, Password: pw, Email: req.Email}
	if err := db.Create(&u).Error; err != nil {
		jsonError(c, http.StatusBadRequest, "username or email already exists")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "username": u.Username})
}

func loginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	var u User
	if err := db.Where("username = ?", req.Username).First(&u).Error; err != nil {
		jsonError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !checkPassword(u.Password, req.Password) {
		jsonError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	tok, err := issueToken(u)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "issue token failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tok})
}

// 处理器：文章 CRUD
func createPostHandler(c *gin.Context) {
	var req PostCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := c.Get("uid")
	p := Post{Title: req.Title, Content: req.Content, UserID: uid.(uint)}
	if err := db.Create(&p).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "create post failed")
		return
	}
	c.JSON(http.StatusCreated, p)
}

func listPostsHandler(c *gin.Context) {
	var posts []Post
	if err := db.Preload("User").Order("id desc").Find(&posts).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "list posts failed")
		return
	}
	c.JSON(http.StatusOK, posts)
}

func getPostHandler(c *gin.Context) {
	var p Post
	id := c.Param("id")
	if err := db.Preload("User").Preload("Comments").First(&p, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, "post not found")
		return
	}
	c.JSON(http.StatusOK, p)
}

func updatePostHandler(c *gin.Context) {
	var req PostUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")
	var p Post
	if err := db.First(&p, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, "post not found")
		return
	}
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if len(updates) == 0 {
		jsonError(c, http.StatusBadRequest, "nothing to update")
		return
	}
	if err := db.Model(&p).Updates(updates).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "update failed")
		return
	}
	c.JSON(http.StatusOK, p)
}

func deletePostHandler(c *gin.Context) {
	id := c.Param("id")
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 先删评论（外键无级联时）
		if err := tx.Where("post_id = ?", id).Delete(&Comment{}).Error; err != nil {
			return err
		}
		// 再删文章
		if err := tx.Delete(&Post{}, id).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		jsonError(c, http.StatusInternalServerError, "delete failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

// 处理器：评论
func createCommentHandler(c *gin.Context) {
	var req CommentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := c.Get("uid")
	pid := c.Param("id")
	// 确保文章存在并加锁读取，防止并发删除
	var post Post
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&post, pid).Error; err != nil {
		jsonError(c, http.StatusNotFound, "post not found")
		return
	}
	cm := Comment{Content: req.Content, UserID: uid.(uint), PostID: post.ID}
	if err := db.Create(&cm).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "create comment failed")
		return
	}
	c.JSON(http.StatusCreated, cm)
}

func listCommentsHandler(c *gin.Context) {
	pid := c.Param("id")
	var cms []Comment
	if err := db.Where("post_id = ?", pid).Order("id asc").Find(&cms).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "list comments failed")
		return
	}
	c.JSON(http.StatusOK, cms)
}

// 数据库初始化
func initDB() (*gorm.DB, error) {
	// 优先读取环境变量 DB_DIALECT 和 DSN
	dialect := strings.ToLower(getEnv("DB_DIALECT", "mysql"))
	dsn := getEnv("MYSQL_DSN", "")

	switch dialect {
	case "mysql":
		if dsn == "" {
			// mysql DSN
			dsn = "root:123456@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
		}
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		file := getEnv("SQLITE_FILE", "blog.db")
		return gorm.Open(sqlite.Open(file), &gorm.Config{})
	default:
		return nil, errors.New("unsupported DB_DIALECT (use mysql or sqlite)")
	}
}

// 路由与启动
func main() {
	var err error
	//数据库启动
	db, err = initDB()
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	//创建表
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	//Gin开始
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 路由健康检查
	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	api := r.Group("/api")
	{
		//注册路由
		api.POST("/register", registerHandler)
		//登录路由
		api.POST("/login", loginHandler)

		// 文章公开读取
		api.GET("/posts", listPostsHandler)
		api.GET("/posts/:id", getPostHandler)

		// 需要登录的接口
		auth := api.Group("")
		//token信息
		auth.Use(authMiddleware())
		{
			// 创建文章
			auth.POST("/posts", createPostHandler)
			// 更新/删除仅作者本人
			auth.PUT("/posts/:id", authorOnlyMiddleware(), updatePostHandler)
			auth.DELETE("/posts/:id", authorOnlyMiddleware(), deletePostHandler)

			// 创建评论
			auth.POST("/posts/:id/comments", createCommentHandler)
		}
		//评论list信息
		api.GET("/posts/:id/comments", listCommentsHandler)
	}

	//web访问初始化
	addr := getEnv("ADDR", ":8080")
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}

}

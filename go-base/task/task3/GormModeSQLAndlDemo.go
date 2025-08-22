package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 表
type User struct {
	//gorm.Model
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"uniqueIndex;size:100"`
	// 一对多关系: 一个用户可以有多篇文章
	Posts []Post `gorm:"foreignKey:UserID"`
	// 文章数量统计
	PostCount int
}

// Post 表
type Post struct {
	//gorm.Model
	ID      int    `gorm:"primaryKey"`
	Title   string `gorm:"size:200;not null"`
	Content string `gorm:"type:text"`
	// 外键，关联到 User
	UserID uint
	// 一对多关系: 一篇文章可以有多个评论
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentStatus string
}

// Comment 表
type Comment struct {
	//gorm.Model
	ID      int    `gorm:"primaryKey"`
	Content string `gorm:"type:text"`
	// 外键，关联到 Post
	PostID uint
}

// AfterCreate Post 的钩子函数，创建后更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	var postCount int64
	if err := tx.Model(&Post{}).Where("user_id = ?", p.UserID).Count(&postCount).Error; err != nil {
		return err
	}

	if err := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", postCount).Error; err != nil {
		return err
	}

	return nil
}

// AfterDelete Comment 的钩子函数，删除后检查文章的评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	if commentCount == 0 {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论").Error; err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// 数据库连接
	dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建 User、Post、Comment 三个表
	//err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	//if err != nil {
	//	log.Fatal("创建失败:", err)
	//}
	//log.Println("数据库表创建成功！")

	//插入信息
	// 插入用户信息
	//var userA, userB User
	//db.FirstOrCreate(&userA, User{ID: 1, Name: "A用户"})
	//db.FirstOrCreate(&userB, User{ID: 2, Name: "B用户"})
	//// 插入文章信息
	//var postA, postB, postC Post
	//db.FirstOrCreate(&postA, Post{ID: 1, Title: "GoLang Base", Content: "GOGOGOGOGOGOGO", UserID: 1})
	//db.FirstOrCreate(&postB, Post{ID: 2, Title: "Solidity Base", Content: "SOSOSOSOOSOSOSO", UserID: 1})
	//db.FirstOrCreate(&postC, Post{ID: 3, Title: "React Base", Content: "RERERERERERERERERE", UserID: 2})
	//// 插入评论信息
	//var commentA, commentB, commentC, commentD, commentE Comment
	//db.FirstOrCreate(&commentA, Comment{ID: 1, Content: "好", PostID: 1})
	//db.FirstOrCreate(&commentB, Comment{ID: 2, Content: "好", PostID: 1})
	//db.FirstOrCreate(&commentC, Comment{ID: 3, Content: "好", PostID: 1})
	//db.FirstOrCreate(&commentD, Comment{ID: 4, Content: "好", PostID: 2})
	//db.FirstOrCreate(&commentE, Comment{ID: 5, Content: "好", PostID: 3})

	// 更新文章的评论状态
	db.Model(&Post{}).Where("id IN (?)", []int{1, 2, 3}).Update("comment_status", "有评论")

	//查询用户的文章和评论
	var user User
	err = db.Preload("Posts.Comments").First(&user, "id = ?", 1).Error
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("用户: %s\n", user.Name)
	for _, post := range user.Posts {
		fmt.Printf("文章: %s\n\n", post.Title)
		for _, comment := range post.Comments {
			fmt.Printf("评论: %s\n", comment.Content)
		}
	}

	//查询评论数量做多的文章
	var topOne Post
	err = db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as commentCount").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("commentCount DESC").
		Limit(1).
		Find(&topOne).Error
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("评论最多的文章: %s\n", topOne.Title)

	//删除评论
	//db.Delete(&commentA)
	//db.Delete(&commentB)

	// 检查文章1的评论状态
	var post1 Post
	db.First(&post1, 1)
	fmt.Printf("文章1的评论状态: %s\n", post1.CommentStatus)

}

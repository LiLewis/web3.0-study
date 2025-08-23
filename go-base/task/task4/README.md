# 个人博客系统的后端（Gin + GORM + JWT）


## 运行环境
- Go 1.25.0
- MySQL 8.0.43


## 依赖安装（go.mod）
url:
https://github.com/LiLewis/web3.0-study/tree/main/go-base/task/task4/go.mod


# 初始化依赖
$ go mod tidy


## 配置 & 启动
# 方式一：MySQL（默认）
export DB_DIALECT=mysql
export MYSQL_DSN="root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
export JWT_SECRET="change_me"
export ADDR=":8080"


go run main.go


# 方式二：SQLite
无


## API 接口
- POST /api/register {username,password,email}
- POST /api/login {username,password} → {token}
- GET /api/posts
- GET /api/posts/:id
- POST /api/posts (Bearer <token>) {title,content}
- PUT /api/posts/:id (Bearer <token>, 作者)
- DELETE /api/posts/:id (Bearer <token>, 作者)
- POST /api/posts/:id/comments (Bearer <token>) {content}
- GET /api/posts/:id/comments


## Postman 调用
- 注册 → 登录 → 在 Authorization 里选择 Bearer Token 粘贴 token
- 先创建文章，再创建评论
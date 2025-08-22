package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Book 结构体
type Book struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Price  int64  `db:"price"`
}

func main() {
	//dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("数据库连接失败: %v", err)
	//}
	////创建表
	//db.AutoMigrate(&Book{})
	//
	////插入书籍信息
	//var bookA, bookB, bookC, bookD Book
	//db.FirstOrCreate(&bookA, Book{ID: 1, Title: "西游记", Author: "吴承恩", Price: 100})
	//db.FirstOrCreate(&bookB, Book{ID: 2, Title: "水浒传", Author: "施耐庵", Price: 90})
	//db.FirstOrCreate(&bookC, Book{ID: 3, Title: "三国演义", Author: "罗贯中", Price: 80})
	//db.FirstOrCreate(&bookD, Book{ID: 4, Title: "红楼梦", Author: "曹雪芹", Price: 70})

	// 数据库连接 DSN
	dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("数据库连接失败:", err)
	}
	defer db.Close()

	// 定义存放结果的切片
	var books []Book
	// 查询价格大于 50 的书籍
	query := `SELECT id, title, author, price FROM books WHERE price > ?`
	err = db.Select(&books, query, 50)
	if err != nil {
		log.Fatalln("查询失败:", err)
	}

	// 打印结果
	for _, b := range books {
		fmt.Printf("ID:%d, Title:%s, Author:%s, Price:%d\n",
			b.ID, b.Title, b.Author, b.Price)
	}
}

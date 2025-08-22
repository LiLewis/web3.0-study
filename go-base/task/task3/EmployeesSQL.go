package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func main() {
	//dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("数据库连接失败: %v", err)
	//}
	////创建表
	//db.AutoMigrate(&Employee{})
	//
	////插入员工信息
	//var employeeA, employeeB Employee
	//db.FirstOrCreate(&employeeA, Employee{ID: 1, Name: "Lewis", Department: "技术部", Salary: 1000})
	//db.FirstOrCreate(&employeeB, Employee{ID: 2, Name: "Leo", Department: "技术部", Salary: 800})

	dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("数据库连接失败:", err)
	}
	defer db.Close()

	// 1. 查询所有技术部员工
	var employees []Employee
	err = db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	fmt.Println("技术部员工:", employees)

	// 2. 查询工资最高的员工
	var topEmployee Employee
	err = db.Get(&topEmployee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	fmt.Println("工资最高的员工:", topEmployee)
}

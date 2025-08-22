package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Student 表
type Student struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Age   int
	Grade string
}

func main() {
	//数据库连接
	dsn := "root:XXXXXX@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	//创建Student 表
	db.AutoMigrate(&Student{})

	//插一条数据
	student := Student{Name: "Lewis", Age: 18, Grade: "家里蹲大学"}
	if err := db.Create(&student).Error; err != nil {
		log.Fatalf("插入失败: %v", err)
	}
	fmt.Println("插入成功:", student)

	//大于16岁的学生
	var students []Student
	if err := db.Where("age > ?", 16).Find(&students).Error; err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Println("查询结果:", students)

	//更改年级信息
	if err := db.Model(&Student{}).Where("name = ?", "Lewis").Update("grade", "天才幼稚园").Error; err != nil {
		log.Fatalf("更新失败: %v", err)
	}
	fmt.Println("更新成功")

	//删除小于15的学生信息
	if err := db.Where("age > ?", 15).Delete(&Student{}).Error; err != nil {
		log.Fatalf("删除失败: %v", err)
	}
	fmt.Println("删除成功")

}

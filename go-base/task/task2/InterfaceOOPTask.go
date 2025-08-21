package main

import (
	"fmt"
	"math"
)

// Shape interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 和 Circle 结构体
type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

// Area 求矩形面积
func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Perimeter 求矩形边长
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Area 求圆形面积
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter 求圆形周长
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Person 实体类
type Person struct {
	Name string
	Age  int
}

// EmployeeVo VO
type EmployeeVo struct {
	Person
	EmployeeID string
}

func (e EmployeeVo) EmployeeInfo() {
	fmt.Printf("姓名: %s, 年龄: %d, 工号: %s\n", e.Name, e.Age, e.EmployeeID)
}

func main() {
	fmt.Println("------------Task1 Start------------")
	//矩形入参
	rect := Rectangle{Width: 10, Height: 5}
	//圆入参
	circle := Circle{Radius: 5}

	//接口变量接收
	var s Shape

	//shapes = append(shapes, &rectangle)
	s = rect
	fmt.Printf("矩形: Area=%.2f, Perimeter=%.2f\n", s.Area(), s.Perimeter())
	//shapes = append(s, &circle)
	s = circle
	fmt.Printf("圆形: Area=%.2f, Perimeter=%.2f\n", s.Area(), s.Perimeter())
	fmt.Println("------------Task1 End------------")

	fmt.Println("------------Task2 Start------------")
	emp := EmployeeVo{
		Person: Person{
			Name: "Lewis",
			Age:  18,
		},
		EmployeeID: "NB-007",
	}
	emp.EmployeeInfo()
	fmt.Println("------------Task2 End------------")
}

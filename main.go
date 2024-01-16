package main

import (
	"fmt"

	"github.com/gin-gonic/gin" // Gin Web Framework
	
	"gorm.io/gorm"				// GORM ORM
  	"gorm.io/driver/sqlite"	// SQLite Driver

)

// MVC = Model View Controller
// User model
type User struct {
	gorm.Model		// ID, CreatedAt, UpdatedAt, DeletedAt
	Name string		// Name field
}

var db *gorm.DB		// Database

func main() {
	// err คือ ตัวแปรที่ใช้เก็บ error ที่เกิดขึ้น
	var err error

	// สร้าง database และเก็บไว้ในตัวแปร db
	db, err := gorm.Open(sqlite.Open("./mydb.sqlite"), &gorm.Config{})
	//db, err = gorm.Open("sqlite3", "./mydb.sqlite")

	// ถ้าเกิด error ให้ panic และแสดงข้อความ "failed to connect database
	if err != nil {
		panic("failed to connect database")
	}
	// ถ้าสิ้นสุดการทำงานของฟังก์ชัน main ให้ปิดการเชื่อมต่อ database
	//defer db.Close()

	// สร้าง table ใน database โดยใช้ struct User
	db.AutoMigrate(&User{})

	// สร้าง router โดยใช้ gin.Default() ซึ่งเป็นการสร้าง router ที่มี middleware ต่างๆ อยู่แล้ว
	r := gin.Default()

	// สร้าง route สำหรับเรียกใช้งาน API
	r.GET("/users", GetUsers)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	// print Start server
	fmt.Println("Start server...")

	// รัน server ที่ port 5000
	r.Run(":5000")
}

// ฟังก์ชัน GetUsers ใช้สำหรับเรียกข้อมูล user ทั้งหมด
func GetUsers(c *gin.Context) {
	// สร้างตัวแปร users และเก็บข้อมูล user ทั้งหมดลงในตัวแปรนั้น
	var users []User

	// ถ้าเกิด error ให้แสดงข้อความ "Error retrieving users"
	if err := db.Find(&users).Error; err != nil {
		// แสดงข้อความ "Error retrieving users" และ HTTP status code 500
		c.JSON(500, gin.H{"error": "Error retrieving users"})
		return
	}
	c.JSON(200, users)	// 200 คือ HTTP status ที่บอกว่าสำเร็จ
}

// ฟังก์ชัน CreateUser ใช้สำหรับสร้าง user
func CreateUser(c *gin.Context) {
	// สร้างตัวแปร user และเก็บข้อมูล user ที่ส่งมาจาก client ลงในตัวแปรนั้น
	var user User

	// BindJSON คือ ฟังก์ชันที่ใช้สำหรับแปลงข้อมูลที่ส่งมาจาก client ให้เป็น JSON
	c.BindJSON(&user)

	// ถ้าเกิด error ให้แสดงข้อความ "Error creating user"
	// Create คือ ฟังก์ชันที่ใช้สำหรับสร้างข้อมูลใน database
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(200, user) // 200 คือ HTTP status ที่บอกว่าสำเร็จ
}

// ฟังก์ชัน UpdateUser ใช้สำหรับแก้ไข user
func UpdateUser(c *gin.Context) {

	// สร้างตัวแปร id และเก็บค่า id ที่ส่งมาจาก client ลงในตัวแปรนั้น
	id := c.Param("id")

	// สร้างตัวแปร user และเก็บข้อมูล user ที่ส่งมาจาก client ลงในตัวแปรนั้น
	var user User

	// ค้นหา user ที่มี id ตามที่ส่งมาจาก client
	db.First(&user, id)

	// ถ้าไม่พบ user ให้แสดงข้อความ "User not found" และ HTTP status code 404
	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// BindJSON คือ ฟังก์ชันที่ใช้สำหรับแปลงข้อมูลที่ส่งมาจาก client ให้เป็น JSON
	c.BindJSON(&user)

	// Save คือ ฟังก์ชันที่ใช้สำหรับบันทึกข้อมูลใน database
	db.Save(&user)

	c.JSON(200, user) // 200 คือ HTTP status ที่บอกว่าสำเร็จ
}

// ฟังก์ชัน DeleteUser ใช้สำหรับลบ user
func DeleteUser(c *gin.Context) {
	// สร้างตัวแปร id และเก็บค่า id ที่ส่งมาจาก client ลงในตัวแปรนั้น
	id := c.Param("id")

	// สร้างตัวแปร user และเก็บข้อมูล user ที่มี id ตามที่ส่งมาจาก client ลงในตัวแปรนั้น
	var user User

	// ค้นหา user ที่มี id ตามที่ส่งมาจาก client
	db.First(&user, id)

	// ถ้าไม่พบ user ให้แสดงข้อความ "User not found" และ HTTP status code 404
	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Delete คือ ฟังก์ชันที่ใช้สำหรับลบข้อมูลใน database
	db.Delete(&user)

	// แสดงข้อความ "User deleted"
	c.JSON(200, gin.H{"success": "User deleted"})
}

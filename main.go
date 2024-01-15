package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name string
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "./mydb.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	r := gin.Default()

	r.GET("/users", GetUsers)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	r.Run(":5000")
}

func GetUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error retrieving users"})
		return
	}
	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	db.First(&user, id)
	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.BindJSON(&user)
	db.Save(&user)

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	db.First(&user, id)
	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	db.Delete(&user)

	c.JSON(200, gin.H{"success": "User deleted"})
}

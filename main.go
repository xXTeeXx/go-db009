package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
func connectToMariaDB() (*gorm.DB, error) {
    dsn := "root:11111@tcp(localhost:3306)/usergo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique"`
    Email    string
}
func main() {
    db, err := connectToMariaDB()
    if err != nil {
        log.Fatal(err)
    }
    // defer db.Close()

    // Perform database migration
    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal(err)
    }

    // Your CRUD operations go here
}
func createUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
func getUserByID(db *gorm.DB, userID uint) (*User, error) {
    var user User
    result := db.First(&user, userID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
func updateUser(db *gorm.DB, user *User) error {
    result := db.Save(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
func deleteUser(db *gorm.DB, user *User) error {
    result := db.Delete(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

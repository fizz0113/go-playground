package main

import (
	"flag"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	username = flag.String("username", "admin", "Specify username")
	password = flag.String("password", "admin", "Specify password")
	hostName = flag.String("hostName", "mysql", "Specify hostName")
	hostPort = flag.Uint("hostPort", 3306, "Specify hostPort")
	database = flag.String("database", "defaultdb", "Specify database")
)

type User struct {
	gorm.Model
	Name     string
	Account  string
	Password string
}

func main() {
	flag.Parse()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		*username,
		*password,
		*hostName,
		*hostPort,
		*database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("db connected")

	models := []interface{}{
		&User{},
	}
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
		defer db.Migrator().DropTable(model)
	}

	// Create User
	user := User{Name: "ian", Account: "ian", Password: "123123"}
	err = db.Create(&user).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", user)

	// Get User
	findUser := User{}
	err = db.Where("id = ?", user.ID).Find(&user).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("findUser:", findUser)

	// Update User
	updateUser := User{Name: "ian2", Account: "ian2", Password: "321321"}
	err = db.Save(&updateUser).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("updateUser:", updateUser)

	// Find N Users
	findNUsers := []User{}
	n := 5
	err = db.Limit(n).Order("created_at desc").Find(&findNUsers).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("findNUsers:", findNUsers)

	// Delete User
	err = db.Where("id = ?", user.ID).Delete(&User{}).Error
	if err != nil {
		panic(err)
	}
}

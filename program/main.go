package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Name string
}

func main() {
	connectionString := os.Getenv("CONNECTION_STRING")
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&User{})
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":80"
	}
	e := echo.New()
	e.GET("/", hello)
	e.GET("/:user", getOneUser)
	e.GET("/allUser", getUser)
	e.POST("/post", postName)
	if err := e.Start(port); err != nil {
		fmt.Println(err)
	}
}

func hello(c echo.Context) error {
	return c.String(200, "hello gais!!")
}

func postName(c echo.Context) error {
	name := User{}
	c.Bind(&name)
	if err := DB.Save(&name).Error; err != nil {
		return c.String(500, err.Error())
	}
	return c.String(200, fmt.Sprintf("Hello %s", name.Name))
}

func getUser(c echo.Context) error {
	var users []User = []User{}
	if err := DB.Find(&users).Error; err != nil {
		fmt.Println(err)
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func getOneUser(c echo.Context) error {
	userName := c.Param("user")
	var users User
	if err := DB.Find(&users, "Name=?", userName).Error; err != nil {
		fmt.Println(err)
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

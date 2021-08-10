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
	DB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
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
	e.GET("/:name", helloName)
	e.GET("/user", getUser)
	if err := e.Start(port); err != nil {
		fmt.Println(err)
	}
}

func hello(c echo.Context) error {
	return c.String(200, "Hajimemashite")
}

func helloName(c echo.Context) error {
	name := c.Param("name")
	return c.String(200, fmt.Sprintf("Hello %s", name))
}

func getUser(c echo.Context) error {
	var users []User
	DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}

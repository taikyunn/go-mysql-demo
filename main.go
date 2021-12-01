package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	USER := os.Getenv("API_USER")
	PASS := os.Getenv("API_PASS")
	ADDRESS := os.Getenv("API_ADDRESS")
	DBMS := "mysql"
	DBNAME := "go_mysql_demo"

	// 本番環境の場合は値を上書きする
	if os.Getenv("DB_ENV") == "production" {
		USER = os.Getenv("DB_USER")
		PASS = os.Getenv("DB_PASS")
		ADDRESS = os.Getenv("DB_ADDRESS")
	}

	CONNECT := USER + ":" + PASS + "@tcp(" + ADDRESS + ":3306)" + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	db := gormConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	r.GET("/", func(c *gin.Context) {
		db := gormConnect()
		var user []User
		db.Order("created_at asc").Find(&user)
		c.HTML(200, "index.html", gin.H{
			"user": user,
		})
	})

	r.POST("/new", func(c *gin.Context) {
		db := gormConnect()
		name := c.PostForm("name")
		email := c.PostForm("email")
		db.Create(&User{Name: name, Email: email})
		defer db.Close()

		c.Redirect(302, "/")
	})

	// 本番環境でもdevelopでも動くようにする
	// API_URLが存在する：本番環境 存在しない：developとなる
	// if urlENV := os.Getenv("API_URL"); urlENV != "" {
	// 	port := os.Getenv("PORT")
	// 	r.Run(":" + port)
	// }

	r.Run(":3000")
}

// Package main provides ...
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Songmu/go-httpdate"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type TransFee struct {
	gorm.Model
	Fee    int
	UserId uint
	Date   time.Time
}

type userName struct {
	gorm.Model
	//ID   int
	Name string
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "3chiku.sqlite3")

	// エラー処理
	if err != nil {
		panic("database can not be opened ! (dbInit)")
	}
	db.AutoMigrate(&TransFee{})
	db.AutoMigrate(&userName{})
	defer db.Close()
}

// データベースにメンバーを追加
func dbUserInsert(name string) {
	db, err := gorm.Open("sqlite3", "3chiku.sqlite3")

	// 例外処理
	if err != nil {
		panic("Data base could not be opened (dbUserInsert)")
	}
	db.Create(&userName{Name: name})
	defer db.Close()
}

func dbSelectID(name string) uint {
	db, err := gorm.Open("sqlite3", "3chiku.sqlite3")

	// 例外処理
	if err != nil {
		panic("Database could not be opened (dbSelectID)")
	}
	var username userName
	db.Where("name = ?", name).First(&username)
	db.Close()
	return username.ID
}

// データベースに交通費とuseridを登録
func dbFeeInsert(fee int, userid uint, date time.Time) {
	db, err := gorm.Open("sqlite3", "3chiku.sqlite3")

	// 例外処理
	if err != nil {
		panic("Data base could not be opened (dbUserInsert)")
	}
	db.Create(&TransFee{Fee: fee, UserId: userid, Date: date})
	defer db.Close()
}
func userdbGetAll() []userName {
	db, err := gorm.Open("sqlite3", "3chiku.sqlite3")
	if err != nil {
		panic("Database could not be opened ! (userdbGetAll)")
	}

	var username []userName
	db.Order("created_at desc").Find(&username)
	db.Close()
	return username
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// initialize database
	dbInit()

	// main page
	router.GET("/", func(ctx *gin.Context) {
		// name := userdbGetAll()

		ctx.HTML(200, "test.html", gin.H{})
	})

	// useradd page get
	router.GET("/useradd", func(ctx *gin.Context) {
		// name := userdbGetAll()

		ctx.HTML(200, "useradd.html", gin.H{
			// "name": name,
		})
	})
	router.POST("/useradd/new", func(ctx *gin.Context) {
		name := ctx.PostForm("name")

		// Insert data to database
		dbUserInsert(name)
		// go to Confirmation page
		ctx.HTML(200, "useraddconfirm.html", gin.H{
			"name": name,
		})

		// ctx.Redirect(302, "/")

	})

	// enterfee page get
	router.GET("/enterfee", func(ctx *gin.Context) {
		name := userdbGetAll()

		ctx.HTML(200, "enterfee.html", gin.H{
			"name": name,
		})
	})
	router.POST("/enterfee/new", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		fmt.Println(ctx.PostForm("date"))
		date, _ := httpdate.Str2Time(ctx.PostForm("date"), nil)
		fmt.Println(date)
		fee, _ := strconv.Atoi(ctx.PostForm("fee"))
		userid := dbSelectID(name)

		dbFeeInsert(fee, userid, date)
		ctx.HTML(200, "useraddconfirm.html", gin.H{
			"name": name,
			"date": date,
			"fee":  fee,
		})
	})
	router.StaticFile("/self.png", "./Static/self.png")
	router.GET("/taro", func(ctx *gin.Context) {
		ctx.HTML(200, "taro.html", gin.H{})
	})
	router.Run()
}

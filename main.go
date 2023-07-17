package main

import (
	"fmt"
	"learn-gome/auth"
	"learn-gome/middleware"
	"log"
	"net/http"
	"os"

	// package used to read the .env file
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq" // postgres golang driver
)

type newStudent struct {
	Student_id       uint64 `json:"student_id"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

//insert data
func postHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent

	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusOK, gin.H{
		"message": "success created",
		"data":    newStudent,
	})

}

//getall data
func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []newStudent //pakai list karna hasil responnya <1 row

	db.Find(&newStudent)
	c.JSON(http.StatusOK, gin.H{
		"message": "success get all data",
		"data":    newStudent,
	})
}

//get data by id
func getHandler(c *gin.Context, db *gorm.DB) {

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success find by id",
		"data":    newStudent,
	})
}

func putHandler(c *gin.Context, db *gorm.DB) {

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
		})
		return
	}

	var reqStudent = newStudent

	c.Bind(&reqStudent)

	db.Model(&newStudent).Where("student_id=?", studentId).Update(reqStudent)

	c.JSON(http.StatusOK, gin.H{
		"message": "success updated",
		"data":    reqStudent,
	})

}

func deleteHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent

	studentId := c.Param("student_id")

	db.Delete(&newStudent, "student_id=?", studentId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Deleted",
	})

}

func setupRouter() *gin.Engine {
	//setup connection db
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error load env")
	}

	conn := os.Getenv("POSTGREE_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	r := gin.Default()

	//home
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message" : "success"})
	})

	//endpoint login
	r.POST("/login", auth.LoginHandler)

	//endpoint post
	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})

	//endpoint getall
	r.GET("/student", middleware.AuthValidation, func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})

	//endpoint getone
	r.GET("/student/:student_id", func(ctx *gin.Context) {
		getHandler(ctx, db)
	})

	//endpoint updateone
	r.PUT("/student/:student_id", func(ctx *gin.Context) {
		putHandler(ctx, db)
	})

	//endpoint delete
	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		deleteHandler(ctx, db)
	})

	return r
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&newStudent{})

	data := &newStudent{}
	if db.Find(&data).RecordNotFound() {
		fmt.Print("======= run seeder user =======")
		seederUser(db)
	}
}

//seeder untuk pengisian data awal
func seederUser(db *gorm.DB) {

	data := newStudent{
		Student_id:       1,
		Student_name:     "Aura",
		Student_age:      10,
		Student_address:  "Malang",
		Student_phone_no: "279458392",
	}

	db.Create(&data)
}

func main() {
	r := setupRouter()

	r.Run(":8080")

}
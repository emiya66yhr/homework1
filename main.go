package main

import (
	"fmt"
	_ "homework1/docs" // Swagger 文档自动生成的包
	"homework1/handlers"
	"homework1/models"
	"homework1/utils"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Student Management API
// @version 1.0
// @description This is a sample server for managing students and their scores.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	// 初始化 Gin 路由
	r := gin.Default()
	r.Use(cors.Default())

	// 配置 API 路由
	r.POST("/student", handlers.AddStudent)                 // 添加学生
	r.PUT("/student/:id", handlers.UpdateStudent)           // 更新学生
	r.DELETE("/student/:id", handlers.DeleteStudent)        // 删除学生
	r.GET("/student/:id", handlers.GetStudent)              // 查询学生
	r.POST("/student/:id/score/:course", handlers.AddScore) // 为学生添加成绩

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 加载 CSV 文件数据
	loadCSVData()

	// 启动服务器
	r.Run(":8080")
}

// loadCSVData 从 CSV 文件加载数据
func loadCSVData() {
	var wg sync.WaitGroup
	ch := make(chan *models.Student, 100)
	errCh := make(chan utils.ImportCSVError, 100)

	wg.Add(1)
	go utils.LoadStudentsFromCSV("students.csv", ch, errCh, &wg)

	go func() {
		for student := range ch {
			utils.AddStudentToMemory(student)
		}
	}()

	wg.Wait()
	close(ch)
	close(errCh)

	for err := range errCh {
		fmt.Printf("Error: %v\n", err)
	}
}

package test

import (
	"bytes"
	"encoding/json"
	"homework1/handlers"
	"homework1/models"
	"homework1/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 创建测试引擎
func setupRouter() *gin.Engine {
	r := gin.Default()
	// 路由配置
	r.POST("/student", handlers.AddStudent)                 // 添加学生
	r.PUT("/student/:id", handlers.UpdateStudent)           // 更新学生
	r.DELETE("/student/:id", handlers.DeleteStudent)        // 删除学生
	r.GET("/student/:id", handlers.GetStudent)              // 查询学生
	r.POST("/student/:id/score/:course", handlers.AddScore) // 为学生添加成绩
	return r
}

// 测试添加学生功能
func TestAddStudent(t *testing.T) {
	// 初始化存储数据
	storage.StudentData = make(map[string]*models.Student)

	// 构造请求数据
	newStudent := models.Student{
		ID:     "1020",
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 82, "线性代数": 79, "复变函数": 81, "马克思主义基本原理": 80},
	}
	jsonData, _ := json.Marshal(newStudent)

	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/student", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	r := setupRouter()
	r.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"学生添加成功"}`, w.Body.String())
}

// 测试更新学生信息
func TestUpdateStudent(t *testing.T) {
	// 初始化存储数据
	storage.StudentData = make(map[string]*models.Student)
	storage.StudentData["1020"] = &models.Student{
		ID:     "1020",
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 82, "线性代数": 79, "复变函数": 81, "马克思主义基本原理": 80},
	}

	// 构造更新数据
	updatedStudent := models.Student{
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 85, "线性代数": 82, "复变函数": 88, "马克思主义基本原理": 83},
	}
	jsonData, _ := json.Marshal(updatedStudent)

	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/student/1020", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	r := setupRouter()
	r.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"学生信息更新成功"}`, w.Body.String())
}

// 测试查询学生信息
func TestGetStudent(t *testing.T) {
	// 初始化存储数据
	storage.StudentData = make(map[string]*models.Student)
	storage.StudentData["1020"] = &models.Student{
		ID:     "1020",
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 82, "线性代数": 79, "复变函数": 81, "马克思主义基本原理": 80},
	}

	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/student/1020", nil)

	// 执行请求
	r := setupRouter()
	r.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, 200, w.Code)
	expected := `{
		"id": "1020",
		"name": "郭勇",
		"gender": "男",
		"class": "10B",
		"scores": {
			"微积分": 82,
			"线性代数": 79,
			"复变函数": 81,
			"马克思主义基本原理": 80
		}
	}`
	assert.JSONEq(t, expected, w.Body.String())
}

// 测试删除学生
func TestDeleteStudent(t *testing.T) {
	// 初始化存储数据
	storage.StudentData = make(map[string]*models.Student)
	storage.StudentData["1020"] = &models.Student{
		ID:     "1020",
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 82, "线性代数": 79, "复变函数": 81, "马克思主义基本原理": 80},
	}

	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/student/1020", nil)

	// 执行请求
	r := setupRouter()
	r.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"学生删除成功"}`, w.Body.String())
}

// 测试添加课程成绩
func TestAddScore(t *testing.T) {
	// 初始化存储数据
	storage.StudentData = make(map[string]*models.Student)
	storage.StudentData["1020"] = &models.Student{
		ID:     "1020",
		Name:   "郭勇",
		Gender: "男",
		Class:  "10B",
		Scores: map[string]float64{"微积分": 82, "线性代数": 79, "复变函数": 81, "马克思主义基本原理": 80},
	}

	// 构造请求体
	scoreData := models.ScoreInput{
		Score: 90,
	}
	jsonData, _ := json.Marshal(scoreData)

	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/student/1020/score/微积分", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	r := setupRouter()
	r.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"课程成绩添加成功"}`, w.Body.String())
}

package handlers

import (
	"homework1/models"
	"homework1/storage"
	"homework1/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 是一个通用的 API 响应结构体
// @description 表示一个标准的 API 响应，包含消息和可选的错误信息
type APIResponse struct {
	Message string `json:"message" example:"操作成功"`
	Error   string `json:"error,omitempty" example:"无效的输入"`
}

// AddStudent godoc
// @Summary 添加新学生
// @Description 通过提供学生的详细信息来向系统中添加新学生。
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.Student true "学生信息"
// @Success 200 {object} APIResponse "成功消息"
// @Failure 400 {object} APIResponse "无效输入"
// @Failure 409 {object} APIResponse "学生已存在"
// @Router /student [post]
func AddStudent(c *gin.Context) {
	var newStudent models.Student
	if err := c.ShouldBindJSON(&newStudent); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "请求体格式错误", Message: "请检查传入的字段或格式是否正确"})
		return
	}

	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if _, exists := storage.StudentData[newStudent.ID]; exists {
		c.JSON(http.StatusConflict, APIResponse{Error: "学生已存在", Message: "学号为 " + newStudent.ID + " 的学生已存在"})
		return
	}

	storage.StudentData[newStudent.ID] = &newStudent

	err := utils.UpdateCSVFile("students.csv", storage.StudentData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "服务器内部错误", Message: "更新 CSV 文件失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{Message: "学生添加成功"})
}

// UpdateStudent godoc
// @Summary 更新学生信息
// @Description 通过学生 ID 更新学生的基本信息。
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "学生 ID"
// @Param student body models.Student true "更新后的学生信息"
// @Success 200 {object} APIResponse "成功消息"
// @Failure 400 {object} APIResponse "无效输入"
// @Failure 404 {object} APIResponse "学生未找到"
// @Router /student/{id} [put]
func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var updatedStudent models.Student
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "无效输入", Message: "请检查传入的字段或格式是否正确"})
		return
	}

	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	student, exists := storage.StudentData[id]
	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{Error: "学生未找到", Message: "学号为 " + id + " 的学生未找到"})
		return
	}

	student.Name = updatedStudent.Name
	student.Gender = updatedStudent.Gender
	student.Class = updatedStudent.Class
	student.Scores = updatedStudent.Scores // 更新成绩

	// 更新 CSV 文件
	err := utils.UpdateCSVFile("students.csv", storage.StudentData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "服务器内部错误", Message: "更新 CSV 文件失败，请稍后重试"})
		return
	}
	c.JSON(http.StatusOK, APIResponse{Message: "学生信息更新成功"})
}

// DeleteStudent godoc
// @Summary 删除学生
// @Description 通过学生 ID 从系统中删除学生。
// @Tags students
// @Produce json
// @Param id path string true "学生 ID"
// @Success 200 {object} APIResponse "成功消息"
// @Failure 404 {object} APIResponse "学生未找到"
// @Router /student/{id} [delete]
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if _, exists := storage.StudentData[id]; exists {
		delete(storage.StudentData, id)
		// 更新 CSV 文件
		err := utils.UpdateCSVFile("students.csv", storage.StudentData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{Error: "服务器内部错误", Message: "更新 CSV 文件失败，请稍后重试"})
			return
		}

		c.JSON(http.StatusOK, APIResponse{Message: "学生删除成功"})
	} else {
		c.JSON(http.StatusNotFound, APIResponse{Error: "学生未找到", Message: "学号为 " + id + " 的学生未找到"})
	}
}

// GetStudent godoc
// @Summary 查询学生信息
// @Description 根据学生 ID 获取学生的详细信息。
// @Tags students
// @Produce json
// @Param id path string true "学生 ID"
// @Success 200 {object} models.Student
// @Failure 404 {object} APIResponse "学生未找到"
// @Router /student/{id} [get]
func GetStudent(c *gin.Context) {
	id := c.Param("id")

	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	student, exists := storage.StudentData[id]
	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{Error: "学生未找到", Message: "学号为 " + id + " 的学生未找到"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// AddScore godoc
// @Summary 为学生添加课程成绩
// @Description 为学生的特定课程添加成绩。
// @Tags scores
// @Accept json
// @Produce json
// @Param id path string true "学生 ID"
// @Param course path string true "课程名称"
// @Param score body models.ScoreInput true "课程成绩"
// @Success 200 {object} APIResponse "成功消息"
// @Failure 400 {object} APIResponse "无效输入"
// @Failure 404 {object} APIResponse "学生未找到"
// @Router /student/{id}/score/{course} [post]
func AddScore(c *gin.Context) {
	id := c.Param("id")
	course := c.Param("course")
	var score models.ScoreInput

	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "无效输入", Message: "请检查传入的成绩是否有效"})
		return
	}

	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	student, exists := storage.StudentData[id]
	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{Error: "学生未找到", Message: "学号为 " + id + " 的学生未找到"})
		return
	}

	student.Scores[course] = score.Score
	c.JSON(http.StatusOK, APIResponse{Message: "课程成绩添加成功"})
}

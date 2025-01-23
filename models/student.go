package models

// Student 结构体用于存储学生信息和成绩
type Student struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Gender string             `json:"gender"`
	Class  string             `json:"class"`
	Scores map[string]float64 `json:"scores"` // 存储课程名称和成绩
}

// NewStudent 创建一个新的学生对象
func NewStudent(id, name, gender, className string) *Student {
	return &Student{
		ID:     id,
		Name:   name,
		Gender: gender,
		Class:  className,
		Scores: make(map[string]float64),
	}
}

// ScoreInput 表示为学生添加分数时的输入结构体
type ScoreInput struct {
	Score float64 `json:"score"`
}

// APIResponse 通用的 API 响应结构体
type APIResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"homework1/models"
	"homework1/storage"
	"os"
	"strconv"
	"sync"
)

// ImportCSVError 用于记录 CSV 导入过程中的错误
type ImportCSVError struct {
	Line int   // 错误行号
	Err  error // 错误详情
}

// LoadStudentsFromCSV 从 CSV 文件中加载学生数据并发送到 channel
// filePath 是 CSV 文件路径
// ch 是传递学生数据的 channel
// errCh 是传递错误信息的 channel
// wg 是用来同步 goroutine 的 WaitGroup
func LoadStudentsFromCSV(filePath string, ch chan<- *models.Student, errCh chan<- ImportCSVError, wg *sync.WaitGroup) {
	defer wg.Done()

	// 打开 CSV 文件
	file, err := os.Open(filePath)
	if err != nil {
		errCh <- ImportCSVError{Line: 0, Err: fmt.Errorf("failed to open file: %v", err)}
		return
	}
	defer file.Close()

	// 创建 CSV Reader
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		errCh <- ImportCSVError{Line: 0, Err: fmt.Errorf("failed to read CSV content: %v", err)}
		return
	}

	// 确保文件有标题行
	if len(lines) < 1 {
		errCh <- ImportCSVError{Line: 0, Err: errors.New("empty or invalid CSV file")}
		return
	}
	headers := lines[0] // 获取标题行

	// 遍历每一行数据
	for i, line := range lines {
		// 跳过标题行
		if i == 0 {
			continue
		}

		// 检查基本数据格式
		if len(line) < 4 {
			errCh <- ImportCSVError{Line: i + 1, Err: errors.New("invalid CSV row format, insufficient columns")}
			continue
		}

		// 解析学生基本信息
		id := line[0]
		name := line[1]
		gender := line[2]
		className := line[3]

		// 解析成绩信息
		scores := make(map[string]float64)
		for j := 4; j < len(line); j++ {
			// 标题行的列名作为课程名
			course := headers[j]
			score, err := strconv.ParseFloat(line[j], 64)
			if err != nil {
				continue // 如果解析失败，跳过该成绩
			}
			scores[course] = score
		}

		// 创建学生对象
		student := models.NewStudent(id, name, gender, className)
		student.Scores = scores

		// 将学生对象发送到 channel
		ch <- student
	}
}

// AddStudentToMemory 将学生数据添加到内存中
func AddStudentToMemory(student *models.Student) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	// 检查是否已经存在该学生
	if existingStudent, exists := storage.StudentData[student.ID]; exists {
		// 更新已有学生的信息和成绩
		existingStudent.Name = student.Name
		existingStudent.Gender = student.Gender
		existingStudent.Class = student.Class
		for course, score := range student.Scores {
			existingStudent.Scores[course] = score
		}
	} else {
		// 添加新学生
		storage.StudentData[student.ID] = student
	}
}

func UpdateCSVFile(filePath string, studentData map[string]*models.Student) error {
	// 打开现有 CSV 文件并读取表头
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV content: %v", err)
	}

	// 提取 CSV 表头（第一行）
	if len(lines) < 1 {
		return fmt.Errorf("CSV file is empty or invalid")
	}
	headers := lines[0] // 第一行是表头

	// 创建一个新的 CSV 文件来保存更新后的数据
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入 CSV 表头
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %v", err)
	}

	// 遍历所有学生，写入数据
	for _, student := range studentData {
		record := []string{student.ID, student.Name, student.Gender, student.Class}
		for _, course := range headers[4:] {
			// 如果存在成绩，添加到对应课程
			if score, exists := student.Scores[course]; exists {
				record = append(record, fmt.Sprintf("%.2f", score))
			} else {
				record = append(record, "")
			}
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return nil
}

// ImportFromCSVWithConcurrency 并发导入 CSV 文件并返回错误信息列表
func ImportFromCSVWithConcurrency(filePath string) []ImportCSVError {
	var wg sync.WaitGroup
	ch := make(chan *models.Student, 100) // 用于传递学生数据
	errCh := make(chan ImportCSVError, 100)

	// 启动 Goroutine，处理学生数据并添加到内存
	go func() {
		for student := range ch {
			AddStudentToMemory(student)
		}
	}()

	// 启动 Goroutine，加载 CSV 文件数据
	wg.Add(1)
	go LoadStudentsFromCSV(filePath, ch, errCh, &wg)

	// 等待所有 Goroutines 完成
	wg.Wait()

	// 关闭 Channels
	close(ch)
	close(errCh)

	// 收集所有错误信息
	var errors []ImportCSVError
	for err := range errCh {
		errors = append(errors, err)
	}

	return errors
}

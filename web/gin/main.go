/*
タスク作成（POST /tasks）
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "first task", "message": "hello world"}'

全タスク取得(GET /tasks)
curl -X GET http://localhost:8080/tasks

特定タスク取得（GET /tasks/:id）
curl -X GET http://localhost:8080/tasks/1

タスク更新（PUT /tasks/:id）
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "updated title", "message": "updated message"}'

タスク削除（DELETE /tasks/:id）
curl -X DELETE http://localhost:8080/tasks/1
*/




package main

import (
	"net/http"
	"bufio"
	"os"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)
type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Message string `json:"message"`
}

const taskFile = "task_list"

func main() {
	// Create Gin router
	r := gin.Default()
	r.PUT("/tasks/:id", updateTask)
	r.GET("/tasks", listTasks)
	r.GET("/task/:id", getTask)
	r.DELETE("/task/:id", deleteTask)
	// Routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin!")
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":   id,
			"name": "User " + id,
		})
	})

	r.POST("/users", func(c *gin.Context) {
		var user struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	})
	r.POST("/tasks", createTask)

	// Start server
	r.Run(":8080")
}


func createTask(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tasks := readTasks()
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1  // このロジック注意
	}
	task := Task{ID: newID, Title: input.Title, Message: input.Message}
	f,_ := os.OpenFile(taskFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // ?
	defer f.Close()
	data,_ := json.Marshal(task)
	f.Write(append(data, '\n'))

	c.JSON(http.StatusCreated, task)
}

func listTasks(c *gin.Context) {
	tasks := readTasks()
	c.JSON(http.StatusOK, tasks)
}
func getTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tasks := readTasks()
	for _,t := range tasks {
		if t.ID == id {
			c.JSON(http.StatusOK, t)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not Found"})
}

func updateTask(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    tasks := readTasks()

    var input struct {
        Title   string `json:"title" binding:"required"`
        Message string `json:"message" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updated := false
    for i, t := range tasks {
        if t.ID == id {
            tasks[i].Title = input.Title
            tasks[i].Message = input.Message
            updated = true
            break
        }
    }

    if !updated {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    // ファイルを全書き換え（DELETE と同じ方式）
    f, _ := os.Create(taskFile)
    defer f.Close()

    for _, t := range tasks {
        data, _ := json.Marshal(t)
        f.Write(append(data, '\n'))
    }

    c.JSON(http.StatusOK, gin.H{"updated": id})
}


func deleteTask(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    tasks := readTasks()

    var newList []Task
    var deleted *Task

    for _, t := range tasks {
        if t.ID == id {
            deleted = &t
        } else {
            newList = append(newList, t)
        }
    }

    if deleted == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    f, _ := os.Create(taskFile)
    defer f.Close()

    for _, t := range newList {
        data, _ := json.Marshal(t)
        f.Write(append(data, '\n'))
    }

    c.JSON(http.StatusOK, gin.H{"deleted": deleted})
}


func readTasks() []Task {
	var tasks []Task

	f,err := os.Open(taskFile)
	if err != nil {
		return tasks
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var t Task
		json.Unmarshal(scanner.Bytes(), &t)
		tasks = append(tasks, t)
	}
	return tasks
}


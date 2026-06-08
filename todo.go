package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{}
var nextID = 1

func ListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", todos)
}

func addHandler(c *gin.Context) {
	task := c.PostForm("task")
	if strings.TrimSpace(task) == "" {
		c.String(http.StatusBadRequest, "Task is required")
		return
	}

	newTodo := Todo{
		ID:        nextID,
		Task:      task,
		Completed: false,
	}
	todos = append(todos, newTodo)
	nextID++

	//在 addHandler 里，添加任务后调用 saveTodos(todos)
	err := saveTodos(todos)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func completeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = true
			break
		}
	}

	//在 doneHandler 里，修改任务后调用 saveTodos(todos)
	if err = saveTodos(todos); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func deleteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	if err = saveTodos(todos); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func loadTodos() {
	// 1. 读取 todos.json
	data, err := os.ReadFile("data/todos.json")
	if err != nil {
		// 2. 如果文件不存在，返回空的 todos 和 nextID=1
		if !os.IsNotExist(err) {
			todos = []Todo{}
			nextID = 1
		}
		return
	}
	// 3. 如果存在，解析 JSON 到 todos 切片
	err = json.Unmarshal(data, &todos)
	if err != nil {
		todos = []Todo{}
		nextID = 1
		return
	}
	// 4. 计算出下一个可用的 nextID（遍历 todos，取最大 ID + 1）
	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	nextID = maxID + 1
}

func saveTodos(todos []Todo) error {
	// 确保data目录存在
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}
	// 1. 把 todos 转成 JSON 格式（用 json.MarshalIndent 更美观）
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	// 2. 写入 todos.json 文件
	return os.WriteFile("data/todos.json", data, 0644)
}

func main() {
	loadTodos()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", ListHandler)
	router.GET("/complete/:id", completeHandler)
	router.GET("/delete/:id", deleteHandler)
	router.POST("/add", addHandler)
	router.Run("localhost:8080")
}

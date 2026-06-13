package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

// isValidTask 验证任务内容是否有效
func isValidTask(task string) bool {
	return strings.TrimSpace(task) != ""
}

var db *sql.DB

func ListHandler(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "todos:list"
	// 1. 尝试从Redis 取缓存
	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var todos []Todo
		if err := json.Unmarshal([]byte(cached), &todos); err == nil {
			c.HTML(http.StatusOK, "list.html", todos)
			return
		}
	}
	// 2. 缓存未命中
	rows, err := db.Query("SELECT id, task, completed FROM todos")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var td Todo
		if err := rows.Scan(&td.ID, &td.Task, &td.Completed); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		todos = append(todos, td)
	}
	if err = rows.Err(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 3. 把查询结果存到Redis， 设置30s过期
	jsonData, err := json.Marshal(todos)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to marshal todos")
		return
	}
	rdb.Set(ctx, cacheKey, jsonData, 30*time.Second)

	c.HTML(http.StatusOK, "list.html", todos)
}

func addHandler(c *gin.Context) {
	task := c.PostForm("task")
	if !isValidTask(task) {
		c.String(http.StatusBadRequest, "Task is required")
		return
	}

	id, err := addTodo(Todo{
		Task:      task,
		Completed: false,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to add todo")
		return
	}
	fmt.Printf("Added todo: %d, %s\n", id, task)

	invalidateCache()
	c.Redirect(http.StatusSeeOther, "/")
}

func completeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	td, err := queryTodoByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "Todo not found")
		return
	}
	err = updateTodo(Todo{
		ID:        id,
		Task:      td.Task,
		Completed: true,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update todo")
		return
	}

	invalidateCache()
	c.Redirect(http.StatusSeeOther, "/")
}

func deleteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	err = deleteTodo(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	invalidateCache()
	c.Redirect(http.StatusSeeOther, "/")
}

func main() {
	// 配置并连接 MySQL
	cfg := mysql.NewConfig()
	cfg.User = getEnv("DBUSER", "root")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", getEnv("DBHOST", "127.0.0.1"), getEnv("DBPORT", "3306"))
	cfg.DBName = getEnv("DBNAME", "gin_todo")

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	// 连接 Redis
	initRedis()
	// 设置Gin 路由
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", ListHandler)
	router.GET("/complete/:id", completeHandler)
	router.GET("/delete/:id", deleteHandler)
	router.POST("/add", addHandler)

	fmt.Println("Server starting on http://0.0.0.0:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func queryTodoByID(id int) (Todo, error) {
	// An album to hold data from the returned row.
	var td Todo

	row := db.QueryRow("SELECT * FROM todos WHERE id = ?", id)
	if err := row.Scan(&td.ID, &td.Task, &td.Completed); err != nil {
		if err == sql.ErrNoRows {
			return td, fmt.Errorf("queryTodoByID %d: no such id", id)
		}
		return td, fmt.Errorf("queryTodoByID %d: %v", id, err)
	}
	return td, nil
}

func addTodo(td Todo) (int, error) {
	result, err := db.Exec("INSERT INTO todos (task, completed) VALUES (?, ?)", td.Task, td.Completed)
	if err != nil {
		return 0, fmt.Errorf("addTodo: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addTodo: %v", err)
	}
	return int(id), nil
}

func updateTodo(td Todo) error {
	_, err := db.Exec("UPDATE todos SET task = ?, completed = ? WHERE id = ?", td.Task, td.Completed, td.ID)
	if err != nil {
		return fmt.Errorf("updateTodo: %v", err)
	}
	return nil
}

func deleteTodo(id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleteTodo: %v", err)
	}
	return nil
}

var rdb *redis.Client
var cacheKey = "todos:list"

// getEnv 获取环境变量，如果未设置则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
	})
}

func invalidateCache() {
	ctx := context.Background()
	rdb.Del(ctx, cacheKey)
}

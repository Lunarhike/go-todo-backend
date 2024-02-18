package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type todo struct {
	ID     int  `json:"id"`
	Task  string  `json:"task"`
	Completed bool  `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "DATABASE_URL")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
    }))
	router.GET("/api/todos", getTodos)
	router.POST("/api/todos", createTodo)

	router.Run("localhost:8080")
}

func getTodos(c *gin.Context) {


	rows, err := db.Query("SELECT id, task, completed, created_at FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []todo
	for rows.Next() {
		var a todo
		err := rows.Scan(&a.ID, &a.Task, &a.Completed, &a.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, todos)
	
}

func createTodo(c *gin.Context) {
	var awesomeTodo todo
	if err := c.BindJSON(&awesomeTodo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO todos (task) VALUES ($1)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeTodo.Task); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, awesomeTodo)
}
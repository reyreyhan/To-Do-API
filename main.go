package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:@/go-gin-gonic?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo{})
}

func main() {
	router := gin.Default()

	router.GET("/", hello)

	routeTodo := router.Group("/api/todos")
	{
		routeTodo.POST("/", create)
		routeTodo.GET("/", all)
		routeTodo.GET("/:id", find)
		routeTodo.POST("/:id", update)
		routeTodo.DELETE(":id", delete)
	}

	router.Run()

}

func hello(c *gin.Context) {
	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "Hello World!",
		},
	)
}

type (
	todo struct {
		gorm.Model
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	transformedTodo struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}
)

func create(c *gin.Context) {
	todo := todo{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Completed:   false,
	}

	db.Save(&todo)

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "Todo item created successfully!",
			"data":    todo,
		},
	)
}

func all(c *gin.Context) {
	var todos []todo
	var transformTodos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			},
		)
		return
	}

	for _, item := range todos {
		transformTodos = append(
			transformTodos,
			transformedTodo{
				ID:          item.ID,
				Title:       item.Title,
				Description: item.Description,
				Completed:   item.Completed,
			})
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   transformTodos,
		},
	)
}

func find(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			})
		return
	}

	transformTodo := transformedTodo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   transformTodo,
		})
}

func update(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	fmt.Println(todo.ID)

	if todo.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	db.Model(&todo).Update("description", c.PostForm("description"))
	completed, _ := strconv.ParseBool(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Todo updated successfully!",
			"data":    todo,
		})
}

func delete(c *gin.Context) {
	var todo todo
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}

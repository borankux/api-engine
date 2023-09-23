# API Engine

### Example 
```go
package main

import (
	"fmt"
	controller "github.com/borankux/api-engine/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required"`
	Posts []Post `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int    `json:"user_id" binding:"required"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
}

var db *gorm.DB

func InitDB() {
	opened, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = opened.AutoMigrate(&User{}, &Post{})
	if err != nil {
		panic("failed to migrate database")
	}
	db = opened
	fmt.Println("Database connected")
}
func GetDB() *gorm.DB {
	return db
}

func main() {
	InitDB()
	app := gin.Default()
	app.Use(func(c *gin.Context) {
		c.Set("database", GetDB())
		c.Next()
	})
	controller.RegisterResource[User]("users", app.Group("/api"), &controller.ResourceConfiguration{
		With:        []string{"Posts"},
		WithMethods: []string{"LIST", "GET"},
	})
	app.Run(":8080")
}

```
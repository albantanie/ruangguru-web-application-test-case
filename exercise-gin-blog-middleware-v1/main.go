package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Posts = []Post{
	{ID: 1, Title: "Judul Postingan Pertama", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for _, user := range users {
			if user.Username == username && user.Password == password {
				c.Next() // Proceed with the request
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(authMiddleware())

	r.GET("/posts", func(c *gin.Context) {
		idStr := c.Query("id")
		if idStr == "" {
			c.JSON(http.StatusOK, gin.H{"posts": Posts})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
			return
		}

		for _, p := range Posts {
			if p.ID == id {
				c.JSON(http.StatusOK, gin.H{"post": p})
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Postingan tidak ditemukan"})
	})

	r.POST("/posts", func(c *gin.Context) {
		var p Post
		err := c.ShouldBindJSON(&p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		p.ID = len(Posts) + 1
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		Posts = append(Posts, p)
		c.JSON(http.StatusCreated, gin.H{"message": "Postingan berhasil ditambahkan", "post": p})
	})

	return r
}

func main() {
	r := SetupRouter()
	r.Run(":8080")
}

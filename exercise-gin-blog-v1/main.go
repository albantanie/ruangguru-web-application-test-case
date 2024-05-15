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

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"posts": Posts})
		// TODO: answer here
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
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
		// TODO: answer here
	})

	r.POST("/posts", func(c *gin.Context) {
		var p Post
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid request body"})
			return
		}

		id := len(Posts) + 1
		now := time.Now()
		p.ID = id
		p.CreatedAt = now
		p.UpdatedAt = now
		Posts = append(Posts, p)
		c.JSON(http.StatusCreated, gin.H{"message": "Postingan berhasil ditambahakan",
			"post": p})

		// TODO: answer here
	})

	return r
}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}

package controllers

import (
	"crud-api-golang/initializers"
	"crud-api-golang/models"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// 🟢 CREATE POST
func PostsCreate(c *gin.Context) {
	var body struct {
		Body  string
		Title string
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}

	// 🧹 Clear cache karena ada data baru
	initializers.Redis.Del(initializers.Ctx, "all_posts")

	c.JSON(200, gin.H{"post": post})
}

// 🟢 READ ALL (pakai Redis)
func PostsIndex(c *gin.Context) {
	// Cek dari Redis
	cached, err := initializers.Redis.Get(initializers.Ctx, "all_posts").Result()
	if err == nil && cached != "" {
		// Jika ditemukan di Redis
		var posts []models.Post
		json.Unmarshal([]byte(cached), &posts)
		c.JSON(200, gin.H{
			"source": "redis",
			"posts":  posts,
		})
		return
	}

	// Kalau tidak ada di Redis, ambil dari DB
	var posts []models.Post
	initializers.DB.Find(&posts)

	// Simpan hasil query ke Redis selama 60 detik
	jsonData, _ := json.Marshal(posts)
	initializers.Redis.Set(initializers.Ctx, "all_posts", jsonData, 60*time.Second)

	c.JSON(200, gin.H{
		"source": "database",
		"posts":  posts,
	})
}

// 🟢 READ ONE
func PostsShow(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, gin.H{"post": post})
}

// 🟢 UPDATE POST
func PostsUpdate(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var post models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	// 🧹 Clear cache agar data di Redis diperbarui nanti
	initializers.Redis.Del(initializers.Ctx, "all_posts")

	c.JSON(200, gin.H{"post": post})
}

// 🟢 DELETE POST
func PostsDelete(c *gin.Context) {
	id := c.Param("id")
	initializers.DB.Delete(&models.Post{}, id)

	// 🧹 Clear cache setelah delete
	initializers.Redis.Del(initializers.Ctx, "all_posts")

	c.Status(200)
}

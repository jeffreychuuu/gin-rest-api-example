package services

import (
	"encoding/json"
	"fmt"
	"gin-rest-api-example/client"
	"gin-rest-api-example/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

var c *gin.Context

func FindBooks(c *gin.Context) {
	var books []models.Book
	client.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

func FindBook(c *gin.Context) {
	var book models.Book
	// Try getting from Redis
	bookJson, _ := client.Redis.HGet(c, "Book", c.Param("id")).Result()
	json.Unmarshal([]byte(bookJson), &book)

	// Get model if exist
	if bookJson == "" {
		// Get Book from DB
		if err := client.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}
		// Cache Book in Redis
		bookJson, err := json.Marshal(book)
		err = client.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
		ttl, err := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
		// Set Redis Expire Time
		client.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Redis Insertion Success!")
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func CreateBook(input models.CreateBookInput) models.Result {
	// Create book in DB
	book := models.Book{Title: input.Title, Author: input.Author}
	client.DB.Create(&book)

	// Create book in Redis
	bookJson, _ := json.Marshal(book)
	err := client.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
	if err != nil {
		panic("Cannot create book in Redis")
		recover()
	}
	ttl, err := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
	if err != nil {
		panic("Cannot set Redis TTL")
		recover()
	}
	// Set Redis Expire Time
	client.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
	fmt.Println("Redis Insertion Success!")

	return models.Result{Data: string(bookJson)}
}

func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := client.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input models.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update book in DB
	client.DB.Model(&book).Updates(input)

	// Update Book in Redis
	bookJson, err := json.Marshal(book)
	err = client.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
	if err != nil {
		panic("Cannot Update Book in Redis")
		recover()
	}
	ttl, _ := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
	// Set Redis Expire Time
	client.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
	fmt.Println("Redis Update Success!")

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := client.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Delete Book in DB
	client.DB.Delete(&book)

	// Delete Book in Redis
	err := client.Redis.HDel(c, "Book", strconv.FormatUint(uint64(book.ID), 10)).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Redis Update Success!")

	c.JSON(http.StatusOK, gin.H{"data": true})
}

package controllers

import (
	"encoding/json"
	"fmt"
	"gin-rest-api-example/config"
	"gin-rest-api-example/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @Tags Book
// @Summary Find books
// @Success 200 {object} models.Result Successful Return Value
// @Router /books [get]
func FindBooks(c *gin.Context) {
	var books []models.Book
	config.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// @Tags Book
// @Summary Find a book
// @Param id path int true "id"
// @Success 200 {object} models.Result Successful Return Value
// @Router /books/{id} [get]
func FindBook(c *gin.Context) {
	var book models.Book
	// Try getting from Redis
	bookJson, _ := config.Redis.HGet(c, "Book", c.Param("id")).Result()
	json.Unmarshal([]byte(bookJson), &book)

	// Get model if exist
	if bookJson == "" {
		// Get Book from DB
		if err := config.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}
		// Cache Book in Redis
		bookJson, err := json.Marshal(book)
		err = config.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
		ttl, err := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
		// Set Redis Expire Time
		config.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Redis Insertion Success!")
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Create new book
// @Param createBookInput body models.CreateBookInput true "CreateBookInput"
// @Success 200 object models.Result
// @Failure 409 object models.Result
// @Router /books [post]
func CreateBook(c *gin.Context) {
	// Validate input
	var input models.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book in DB
	book := models.Book{Title: input.Title, Author: input.Author}
	config.DB.Create(&book)

	// Create book in Redis
	bookJson, _ := json.Marshal(book)
	err := config.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
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
	config.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
	fmt.Println("Redis Insertion Success!")

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Update a book
// @Param id path int true "id"
// @Param updateBookInput body models.UpdateBookInput true "UpdateBookInput"
// @Success 200 object models.Result
// @Failure 409 object models.Result
// @Router /books/{id} [patch]
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := config.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
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
	config.DB.Model(&book).Updates(input)

	// Update Book in Redis
	bookJson, err := json.Marshal(book)
	err = config.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
	if err != nil {
		panic("Cannot Update Book in Redis")
		recover()
	}
	ttl, _ := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
	// Set Redis Expire Time
	config.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
	fmt.Println("Redis Update Success!")

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Delete a book
// @Param id path int true "id"
// @Success 200 object models.Result
// @Failure 409 object models.Result
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := config.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Delete Book in DB
	config.DB.Delete(&book)

	// Delete Book in Redis
	err := config.Redis.HDel(c, "Book", strconv.FormatUint(uint64(book.ID), 10)).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Redis Update Success!")

	c.JSON(http.StatusOK, gin.H{"data": true})
}

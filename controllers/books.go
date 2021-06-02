package controllers

import (
	"gin-rest-api-example/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// @Tags Book
// @Summary Find books
// @Success 200 {object} models.Result Successful Return Value
// @Router /books [get]
func FindBooks(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// @Tags Book
// @Summary Find a book
// @Param id path int true "id"
// @Success 200 {object} models.Result Successful Return Value
// @Router /books/{id} [get]
func FindBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Create new book
// @Param createBookInput body controllers.CreateBookInput true "CreateBookInput"
// @Success 200 object models.Result 成功后返回值
// @Failure 409 object models.Result 添加失败
// @Router /books [post]
func CreateBook(c *gin.Context) {
	// Validate input
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Update a book
// @Param id path int true "id"
// @Param updateBookInput body controllers.UpdateBookInput true "UpdateBookInput"
// @Success 200 object models.Result 成功后返回值
// @Failure 409 object models.Result 添加失败
// @Router /books/{id} [patch]
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&book).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// @Tags Book
// @Summary Delete a book
// @Param id path int true "id"
// @Success 200 object models.Result 成功后返回值
// @Failure 409 object models.Result 添加失败
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

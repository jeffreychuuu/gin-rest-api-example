package controllers

import (
	"gin-rest-api-example/models"
	"gin-rest-api-example/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Tags Book
// @Summary Find books
// @Success 200 {object} models.Result Successful Return Value
// @Router /books [get]
func FindBooks(c *gin.Context) {
	services.FindBooks(c)
}

// @Tags Book
// @Summary Find a book
// @Param id path int true "id"
// @Success 200 {object} models.Result Successful Return Value
// @Router /books/{id} [get]
func FindBook(c *gin.Context) {
	services.FindBook(c)
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
	book := services.CreateBook(input)

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
	services.UpdateBook(c)
}

// @Tags Book
// @Summary Delete a book
// @Param id path int true "id"
// @Success 200 object models.Result
// @Failure 409 object models.Result
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	services.DeleteBook(c)
}

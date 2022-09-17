package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"todolist/models"

	"log"

	"github.com/gin-gonic/gin"
)

// schema that can validate the user's input to prevent invalid data
type CreateItemInput struct {
	Title     string     `json:"title" binding:"required"`
	Status    string     `json:"status" binding:"required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateItemInput struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
}

type DataPaging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

// POST /items
func CreateItem(c *gin.Context) {
	var inputItem CreateItemInput
	if err := c.ShouldBindJSON(&inputItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Item
	// preprocess title - trim all spaces
	inputItem.Title = strings.TrimSpace(inputItem.Title)

	if inputItem.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title should not be blank"})
		return
	}

	// do not allow "Finished" status when creating a new task
	inputItem.Status = "Doing" // set to default

	toDoItem := models.ToDoItem{Title: inputItem.Title, Status: inputItem.Status,
		CreatedAt: inputItem.CreatedAt, UpdatedAt: inputItem.UpdatedAt}

	if err := models.DB.Create(&toDoItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toDoItem.Id})

}

// GET /v1/items
func GetItems(c *gin.Context) {

	var paging DataPaging

	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if paging.Page <= 0 {
		paging.Page = 1
	}

	if paging.Limit <= 0 {
		paging.Limit = 10
	}

	offset := (paging.Page - 1) * paging.Limit

	log.Println("paging.Page:", paging.Page)
	log.Println("paging.Limit:", paging.Limit)
	log.Println("offset:", offset)
	log.Println("table name:", models.ToDoItem{}.TableName())

	var toDoItems []models.ToDoItem

	if err := models.DB.Table(models.ToDoItem{}.TableName()).
		Count(&paging.Total).
		Offset(offset).
		Order("id desc").
		Find(&toDoItems).Error; err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toDoItems})
}

func ParseId(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GET /v1/items/:id
func ReadItemById(c *gin.Context) {
	var dataItem models.ToDoItem

	// id, err := strconv.Atoi(c.Param("id"))

	id, err := ParseId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	if err := models.DB.Where("id = ?", id).First(&dataItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dataItem})

}

// PATCH /v1/items/:id
func EditItemById(c *gin.Context) {
	var item models.ToDoItem

	id, err := ParseId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	if err := models.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// validate input
	var inputItem UpdateItemInput
	if err := c.ShouldBindJSON(&inputItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(inputItem)
	log.Println(item)

	inputItem.Title = strings.TrimSpace(inputItem.Title)
	if inputItem.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title should not be blank"})
		return
	}

	if inputItem.Status != "Doing" && inputItem.Status != "Finished" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status should be Doing or Finished"})
		return
	}

	toDoItem := models.ToDoItem{
		Title: inputItem.Title, Status: inputItem.Status,
		CreatedAt: item.CreatedAt, UpdatedAt: &time.Time{}}

	models.DB.Model(&item).Updates(toDoItem)
	c.JSON(http.StatusOK, gin.H{"data": item})

}

// DELETE /v1/items/:id
func DeleteItemById(c *gin.Context) {
	var item models.ToDoItem

	id, err := ParseId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	if err := models.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	models.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"data": true})

}

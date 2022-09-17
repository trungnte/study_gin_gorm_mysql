package main

import (
	"net/http"
	"todolist/controllers"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

func main() {

	// connect database
	models.ConnectDatabase()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Service is running!"})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/items", controllers.CreateItem)           // create item
		v1.GET("/items", controllers.GetItems)              // get list of item
		v1.GET("/items/:id", controllers.ReadItemById)      // get an item by ID
		v1.PATCH("/items/:id", controllers.EditItemById)    // edit an item by ID
		v1.DELETE("/items/:id", controllers.DeleteItemById) // delete an item by ID
	}

	r.Run()
}

package api

import (
	"net/http"
	"project-root/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterBrigadeRoutes(router gin.IRouter) {
	brigade := router.Group("/api/brigades")
	{
		brigade.POST("/", createBrigade)
		brigade.GET("/", getAllBrigades)
		brigade.GET("/:id", getBrigade)
		brigade.PUT("/:id", updateBrigade)
		brigade.DELETE("/:id", deleteBrigade)
	}

	task := router.Group("/api/tasks")
	{
		task.POST("/", createTask)
		task.GET("/:id", getTask)
		task.PUT("/:id", updateTask)
		task.DELETE("/:id", deleteTask)
	}
}

func createBrigade(c *gin.Context) {
	var brigade services.Brigade
	if err := c.ShouldBindJSON(&brigade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := services.CreateBrigadeInDB(brigade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getBrigade(c *gin.Context) {
	brigadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brigade ID"})
		return
	}

	brigade, err := services.GetBrigadeByID(brigadeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brigade not found"})
		return
	}

	c.JSON(http.StatusOK, brigade)
}

func getAllBrigades(c *gin.Context) {
	brigades, err := services.GetAllBrigades()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, brigades)
}

func updateBrigade(c *gin.Context) {
	brigadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brigade ID"})
		return
	}

	var brigade services.Brigade
	if err := c.ShouldBindJSON(&brigade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brigade.ID = brigadeID

	response, err := services.UpdateBrigadeInDB(brigade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func deleteBrigade(c *gin.Context) {
	brigadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brigade ID"})
		return
	}

	err = services.DeleteBrigadeByID(brigadeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
func createTask(c *gin.Context) {
	var task services.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.AssignedAt = time.Now()

	response, err := services.CreateTaskInDB(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := services.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task services.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = taskID

	response, err := services.UpdateTaskInDB(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func deleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = services.DeleteTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

package api

import (
	"fmt"
	"net/http"
	"project-root/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(router gin.IRouter) {
	booking := router.Group("/api/bookings")
	{
		booking.POST("/", createBooking)
		booking.GET("/:id", getBooking) // здесь вызываем getBooking
	}
}

func createBooking(c *gin.Context) {
	var request services.Booking

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим дату в формате "YYYY-MM-DD"
	parsedDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	booking := services.Booking{
		SlotID: request.SlotID,
		Date:   parsedDate.Format("2006-01-02"), // Преобразуем в строку для сохранения в БД
		Status: request.Status,
	}

	response, err := services.CreateBookingInDB(booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getBooking(c *gin.Context) {
	bookingID := c.Param("id")

	// Преобразование bookingID из string в int
	id, err := strconv.Atoi(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	fmt.Printf("Fetching booking with ID: %d\n", id) // Логируем ID

	// Получаем бронирование из базы данных
	booking, err := services.GetBookingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

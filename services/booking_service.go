package services

import (
	"context"
	"fmt"
	"project-root/config"
)

type Booking struct {
	ID     int    `json:"id"`
	SlotID string `json:"slot_id"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

func CreateBookingInDB(booking Booking) (Booking, error) {
	tx, err := config.DB.Begin(context.Background())
	if err != nil {
		return booking, fmt.Errorf("could not start transaction: %v", err)
	}

	defer tx.Rollback(context.Background()) // Откатываем транзакцию при ошибке

	sql := `INSERT INTO bookings (slot_id, date, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(context.Background(), sql, booking.SlotID, booking.Date, booking.Status).Scan(&booking.ID)
	if err != nil {
		return booking, fmt.Errorf("could not insert booking: %v", err)
	}

	notification := Notification{Message: fmt.Sprintf("New booking created: %d", booking.ID)}
	if err := SendNotification(config.RedisClient, notification); err != nil {
		return booking, fmt.Errorf("could not send notification: %v", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(context.Background()); err != nil {
		return booking, fmt.Errorf("could not commit transaction: %v", err)
	}

	return booking, nil
}

func UpdateBookingInDB(booking Booking) (Booking, error) {
	sql := `UPDATE bookings SET slot_id = $1, date = $2, status = $3 WHERE id = $4`
	_, err := config.DB.Exec(context.Background(), sql, booking.SlotID, booking.Date, booking.Status, booking.ID)
	if err != nil {
		return booking, fmt.Errorf("could not update booking: %v", err)
	}

	// Отправка уведомления
	notification := Notification{Message: fmt.Sprintf("Booking updated: %d", booking.ID)}
	if err := SendNotification(config.RedisClient, notification); err != nil {
		return booking, fmt.Errorf("could not send notification: %v", err)
	}

	return booking, nil
}

func GetBookingByID(bookingID int) (Booking, error) {
	var booking Booking
	sql := `SELECT id, slot_id, date, status FROM bookings WHERE id = $1`
	err := config.DB.QueryRow(context.Background(), sql, bookingID).Scan(&booking.ID, &booking.SlotID, &booking.Date, &booking.Status)
	if err != nil {
		fmt.Printf("Error querying booking with ID %d: %v\n", bookingID, err) // Логирование ошибки
		return booking, fmt.Errorf("could not find booking: %v", err)
	}
	return booking, nil
}

package services

import (
	"context"
	"fmt"
	"project-root/config"
	"time"
)

type Brigade struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Task struct {
	ID          int       `json:"id"`
	BrigadeID   int       `json:"brigade_id"`
	Description string    `json:"description"`
	AssignedAt  time.Time `json:"assigned_at"`
	Status      string    `json:"status"`
}

// Создание бригады
func CreateBrigadeInDB(brigade Brigade) (Brigade, error) {
	sql := `INSERT INTO brigades (name, status) VALUES ($1, $2) RETURNING id`
	err := config.DB.QueryRow(context.Background(), sql, brigade.Name, brigade.Status).Scan(&brigade.ID)
	if err != nil {
		return brigade, fmt.Errorf("could not insert brigade: %v", err)
	}
	return brigade, nil
}

// Получение бригады по ID
func GetBrigadeByID(brigadeID int) (Brigade, error) {
	var brigade Brigade
	sql := `SELECT id, name, status FROM brigades WHERE id = $1`
	err := config.DB.QueryRow(context.Background(), sql, brigadeID).Scan(&brigade.ID, &brigade.Name, &brigade.Status)
	if err != nil {
		return brigade, fmt.Errorf("could not find brigade: %v", err)
	}
	return brigade, nil
}

// Обновление бригады
func UpdateBrigadeInDB(brigade Brigade) (Brigade, error) {
	sql := `UPDATE brigades SET name = $1, status = $2 WHERE id = $3`
	_, err := config.DB.Exec(context.Background(), sql, brigade.Name, brigade.Status, brigade.ID)
	if err != nil {
		return brigade, fmt.Errorf("could not update brigade: %v", err)
	}
	return brigade, nil
}

// Удаление бригады по ID
func DeleteBrigadeByID(brigadeID int) error {
	sql := `DELETE FROM brigades WHERE id = $1`
	_, err := config.DB.Exec(context.Background(), sql, brigadeID)
	if err != nil {
		return fmt.Errorf("could not delete brigade: %v", err)
	}
	return nil
}

func CreateTaskInDB(task Task) (Task, error) {
	sql := `INSERT INTO tasks (brigade_id, description, assigned_at, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := config.DB.QueryRow(context.Background(), sql, task.BrigadeID, task.Description, task.AssignedAt, task.Status).Scan(&task.ID)
	if err != nil {
		return task, fmt.Errorf("could not insert task: %v", err)
	}

	// Отправка уведомления
	notification := Notification{Message: fmt.Sprintf("New task created: %d", task.ID)}
	if err := SendNotification(config.RedisClient, notification); err != nil {
		return task, fmt.Errorf("could not send notification: %v", err)
	}

	return task, nil
}

// Получение задачи по ID
func GetTaskByID(taskID int) (Task, error) {
	var task Task
	sql := `SELECT id, brigade_id, description, assigned_at, status FROM tasks WHERE id = $1`
	err := config.DB.QueryRow(context.Background(), sql, taskID).Scan(&task.ID, &task.BrigadeID, &task.Description, &task.AssignedAt, &task.Status)
	if err != nil {
		return task, fmt.Errorf("could not find task: %v", err)
	}
	return task, nil
}

func UpdateTaskInDB(task Task) (Task, error) {
	sql := `UPDATE tasks SET description = $1, status = $2 WHERE id = $3`
	_, err := config.DB.Exec(context.Background(), sql, task.Description, task.Status, task.ID)
	if err != nil {
		return task, fmt.Errorf("could not update task: %v", err)
	}

	// Отправка уведомления
	notification := Notification{Message: fmt.Sprintf("Task updated: %d", task.ID)}
	if err := SendNotification(config.RedisClient, notification); err != nil {
		return task, fmt.Errorf("could not send notification: %v", err)
	}

	return task, nil
}

// Удаление задачи по ID
func DeleteTaskByID(taskID int) error {
	sql := `DELETE FROM tasks WHERE id = $1`
	_, err := config.DB.Exec(context.Background(), sql, taskID)
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}
	return nil
}

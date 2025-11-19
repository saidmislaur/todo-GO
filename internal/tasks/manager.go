package tasks

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	StatusDone         = "done"
	StatusInProcess    = "in_process"
	StatusNotCompleted = "not_completed"
)

var validStatuses = map[string]bool{
	StatusDone:         true,
	StatusInProcess:    true,
	StatusNotCompleted: true,
}

func NewManager(db *sql.DB) *TaskManager {
	return &TaskManager{DB: db}
}

func (tm *TaskManager) AddTask(userId int, text string) (Task, error) {
	if text == "" {
		return Task{}, errors.New("текст задачи не может быть пустым")
	}

	var task Task
	err := tm.DB.QueryRow(`
		INSERT INTO tasks (user_id, text, status)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, text, status
	`, userId, text, StatusInProcess).Scan(
		&task.ID, &task.UserId, &task.Text, &task.Status,
	)

	if err != nil {
		return Task{}, fmt.Errorf("ошибка создания задачи: %w", err)
	}

	return task, nil
}

func (tm *TaskManager) GetTasksByUser(userId int) ([]Task, error) {
	rows, err := tm.DB.Query(`
		SELECT id, user_id, text, status
		FROM tasks
		WHERE user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasksList []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.UserId, &t.Text, &t.Status); err != nil {
			return nil, err
		}
		tasksList = append(tasksList, t)
	}

	return tasksList, nil
}

func (tm *TaskManager) UpdateTask(id int, updated Task) (Task, error) {
	var existing Task
	err := tm.DB.QueryRow(`
		SELECT id, user_id, text, status
		FROM tasks
		WHERE id = $1
	`, id).Scan(&existing.ID, &existing.UserId, &existing.Text, &existing.Status)

	if err == sql.ErrNoRows {
		return Task{}, fmt.Errorf("задача с id %d не найдена", id)
	}
	if err != nil {
		return Task{}, err
	}

	// если пришёл новый текст
	newText := existing.Text
	if updated.Text != "" {
		newText = updated.Text
	}

	// если пришёл новый статус
	newStatus := existing.Status
	if updated.Status != "" {
		if !validStatuses[updated.Status] {
			return Task{}, fmt.Errorf("недопустимый статус: %s", updated.Status)
		}
		newStatus = updated.Status
	}

	// выполняем обновление
	err = tm.DB.QueryRow(`
		UPDATE tasks
		SET text = $1, status = $2
		WHERE id = $3
		RETURNING id, user_id, text, status
	`, newText, newStatus, id).Scan(
		&existing.ID, &existing.UserId, &existing.Text, &existing.Status,
	)
	if err != nil {
		return Task{}, err
	}

	return existing, nil
}

func (tm *TaskManager) DeleteTask(id int) error {
	res, err := tm.DB.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("задача с id %d не найдена", id)
	}

	return nil
}

package tasks

import "database/sql"

type Task struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
	UserId int    `json:"user_id"`
}

type TaskManager struct {
	DB *sql.DB
}

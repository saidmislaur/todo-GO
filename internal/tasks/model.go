package tasks

import "bufio"

type Task struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
	UserId int    `json:"user_id"`
}

type TaskManager struct {
	Tasks    map[int]Task
	Reader   *bufio.Reader
	FilePath string
}

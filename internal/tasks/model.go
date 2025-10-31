package tasks

import "bufio"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type TaskManager struct {
	Tasks    []Task
	Reader   *bufio.Reader
	FilePath string
}

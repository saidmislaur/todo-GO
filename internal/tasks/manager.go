package tasks

import (
	"bufio"
	"fmt"
	"os"
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

func NewManager(filePath string) *TaskManager {
	return &TaskManager{
		Reader:   bufio.NewReader(os.Stdin),
		FilePath: filePath,
	}
}

func (tm *TaskManager) AddTask(text string) (Task, error) {
	newID := 1
	if len(tm.Tasks) > 0 {
		for id := range tm.Tasks {
			if id >= newID {
				newID = id + 1
			}
		}
	}

	newTask := Task{
		ID:     newID,
		Text:   text,
		Status: StatusInProcess,
	}

	tm.Tasks[newID] = newTask

	if err := tm.SaveTasks(); err != nil {
		delete(tm.Tasks, newID) // откат
		return Task{}, fmt.Errorf("ошибка сохранения: %w", err)
	}
	fmt.Println("Задача добавлена ✅")

	return newTask, nil
}

func (tm *TaskManager) UpdateTask(id int, updated Task) (Task, error) {
	task, exists := tm.Tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("задача с ID %d не найдена", id)
	}

	if updated.Text != "" {
		task.Text = updated.Text
	}

	if updated.Status != "" {
		if !validStatuses[updated.Status] {
			return Task{}, fmt.Errorf("недопустимый статус: %s", updated.Status)
		}
		task.Status = updated.Status
	}

	tm.Tasks[id] = task

	if err := tm.SaveTasks(); err != nil {
		return Task{}, fmt.Errorf("ошибка сохранения: %w", err)
	}

	return task, nil
}

func (tm *TaskManager) DeleteTask(id int) error {
	task, exists := tm.Tasks[id]
	if !exists {
		return fmt.Errorf("задача с id %d не найдена", id)
	}

	delete(tm.Tasks, id)

	if err := tm.SaveTasks(); err != nil {
		return fmt.Errorf("ошибка при сохранении после удаления: %w", err)
	}

	fmt.Println("Удалена задача:", task.Text)

	return nil
}

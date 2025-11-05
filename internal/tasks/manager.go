package tasks

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const (
	StatusDone         = "done"
	StatusInProcess    = "in_process"
	StatusNotCompleted = "not_completed"
)

func NewManager(filePath string) *TaskManager {
	return &TaskManager{
		Reader:   bufio.NewReader(os.Stdin),
		FilePath: filePath,
	}
}

func (tm *TaskManager) AddTask(text string) (Task, error) {
	newTask := Task{
		ID:     rand.Intn(1000),
		Text:   text,
		Status: "in_process",
	}

	tm.Tasks = append(tm.Tasks, newTask)
	if err := tm.SaveTasks(); err != nil {
		return Task{}, fmt.Errorf("ошибка при добавлении: %w", err)
	}
	fmt.Println("Задача добавлена ✅")

	return newTask, nil
}

func (tm *TaskManager) UpdateTask(id int, updated Task) (Task, error) {
	for i, task := range tm.Tasks {
		if task.ID == id {
			if updated.Text != "" {
				tm.Tasks[i].Text = updated.Text
			}

			if updated.Status != "" {
				allowedStatuses := []string{StatusDone, StatusInProcess, StatusNotCompleted}
				valid := false
				for _, s := range allowedStatuses {
					if updated.Status == s {
						valid = true
						break
					}
				}
				if !valid {
					return Task{}, fmt.Errorf("недопустимый статус: %s", updated.Status)
				}
				tm.Tasks[i].Status = updated.Status
			}
			if err := tm.SaveTasks(); err != nil {
				return Task{}, fmt.Errorf("ошибка сохранения: %w", err)
			}

			return tm.Tasks[i], nil
		}
	}

	return Task{}, fmt.Errorf("задача с ID %d не найдена", id)
}

func (tm *TaskManager) DeleteTask(id int) error {
	num := -1

	for i, task := range tm.Tasks {
		if task.ID == id {
			num = i
			break
		}
	}

	if num == -1 {
		return fmt.Errorf("задача с id %d не найдена", id)
	}

	deleted := tm.Tasks[num]
	tm.Tasks = append(tm.Tasks[:num], tm.Tasks[num+1:]...)
	if err := tm.SaveTasks(); err != nil {
		return fmt.Errorf("ошибка при сохранении после удаления: %w", err)
	}

	fmt.Println("Удалена задача:", deleted.Text)

	return nil
}

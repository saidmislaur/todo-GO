package tasks

import (
	"encoding/json"
	"os"
)

func (tm *TaskManager) SaveTasks() error {
	var tasks []Task
	for _, task := range tm.Tasks {
		tasks = append(tasks, task)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tm.FilePath, data, 0644)
}

func (tm *TaskManager) LoadTasks() error {
	data, err := os.ReadFile(tm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			tm.Tasks = make(map[int]Task)
			return nil
		}
		return err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	tm.Tasks = make(map[int]Task)
	for _, task := range tasks {
		tm.Tasks[task.ID] = task
	}

	return nil
}

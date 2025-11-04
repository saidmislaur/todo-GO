package tasks

import (
	"encoding/json"
	"fmt"
	"os"
)

func (tm *TaskManager) SaveTasks() error {
	file, err := os.Create(tm.FilePath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	if err := json.NewEncoder(file).Encode(tm.Tasks); err != nil {
		return fmt.Errorf("ошибка записи JSON: %w", err)
	}

	return nil
}

func (tm *TaskManager) LoadTasks() error {
	file, err := os.Open(tm.FilePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&tm.Tasks); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return nil
}

package tasks

import (
	"encoding/json"
	"log"
	"os"
)

func (tm *TaskManager) SaveTasks() {
	file, err := os.Create("tasks.json")
	if err != nil {
		log.Fatalf("%v", file)
	}

	json.NewEncoder(file).Encode(tm.Tasks)
}

func (tm *TaskManager) LoadTasks() {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Файл задач не найден, создаю новый")
		return
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&tm.Tasks)
}

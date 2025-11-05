package main

import (
	"demo/app/internal/tasks"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	manager := tasks.NewManager("tasks.json")
	if err := manager.LoadTasks(); err != nil {
		fmt.Println("ошибка загрузки задач")
	}

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Сontent-Type ", "application/json")
			json.NewEncoder(w).Encode(manager.Tasks)

		case http.MethodPost:
			var newTask tasks.Task
			if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
				http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
				return
			}

			if newTask.Text == "" {
				http.Error(w, "Текст задачи не может быть пустым", http.StatusBadRequest)
				return
			}

			createdTask, err := manager.AddTask(newTask.Text)
			if err != nil {
				http.Error(w, "ошибка при сохранении задачи", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type ", "aplication/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(createdTask)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			var updated tasks.Task
			if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
				http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
				return
			}
			task, err := manager.UpdateTask(id, updated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type ", "application/json")
			json.NewEncoder(w).Encode(task)

		case http.MethodDelete:
			if err := manager.DeleteTask(id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Сервер запущен: http://localhost:5000")
	http.ListenAndServe(":5000", nil)
}

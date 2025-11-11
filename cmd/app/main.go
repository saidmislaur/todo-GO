package main

import (
	"demo/app/internal/tasks"
	"demo/app/internal/users"
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
	userManager := users.NewManager("users.json")
	if err := userManager.LoadUsers(); err != nil {
		fmt.Println("ошибка загрузки пользователей")
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var newUser users.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}
		if err := userManager.Register(newUser.Username, newUser.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var user users.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		token, err := userManager.Login(user.Username, user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := userManager.GetUserIDByToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Сontent-Type ", "application/json")
			json.NewEncoder(w).Encode(manager.GetTasksByUser(userId))

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

			createdTask, err := manager.AddTask(userId, newTask.Text)
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

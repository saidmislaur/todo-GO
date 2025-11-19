package main

import (
	"demo/app/internal/database"
	"demo/app/internal/tasks"
	"demo/app/internal/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	db := database.Connect()
	defer db.Close()

	taskManager := tasks.NewManager(db)
	userManager := users.NewManager(db)

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

		if newUser.Username == "" || newUser.Password == "" {
			http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
			return
		}

		if err := userManager.Register(newUser.Username, newUser.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь зарегистрирован",
		})
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		userId, err := getUserIdFromAuth(r, userManager)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		switch r.Method {

		case http.MethodGet:
			tasksList, err := taskManager.GetTasksByUser(userId)
			if err != nil {
				http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasksList)

		case http.MethodPost:
			var newTask tasks.Task
			if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
				http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
				return
			}

			if newTask.Text == "" {
				http.Error(w, "Текст задачи обязателен", http.StatusBadRequest)
				return
			}

			created, err := taskManager.AddTask(userId, newTask.Text)
			if err != nil {
				http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(created)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		userId, err := getUserIdFromAuth(r, userManager)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

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
				http.Error(w, "Неверный JSON", http.StatusBadRequest)
				return
			}

			task, err := taskManager.UpdateTask(id, updated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			if task.UserId != userId {
				http.Error(w, "Доступ запрещён", http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)

		case http.MethodDelete:
			err := taskManager.DeleteTask(id)
			if err != nil {
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

func getUserIdFromAuth(r *http.Request, um *users.UserManager) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, fmt.Errorf("Требуется авторизация")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return um.GetUserIDByToken(token)
}

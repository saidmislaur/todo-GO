package ui

import (
	"demo/app/internal/tasks"
	"fmt"
)

const (
	MenuShowTasks  = "1"
	MenuAddTask    = "2"
	MenuDeleteTask = "3"
	MenuUpdateTask = "4"
	MenuExit       = "5"
)

func Run(tm *tasks.TaskManager) bool {
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1. Показать задачи")
		fmt.Println("2. Добавить задачу")
		fmt.Println("3. Удалить задачу")
		fmt.Println("4. Изменить задачу")
		fmt.Println("5. Выход")

		var choice string
		fmt.Scan(&choice)

		switch choice {
		case MenuShowTasks:
			tm.ShowTasks()
		case MenuAddTask:
			if err := tm.AddTask(); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case MenuDeleteTask:
			if err := tm.DeleteTask(); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case MenuUpdateTask:
			if err := tm.UpdateTask(); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case MenuExit:
			fmt.Println("Выход из программы...")
			return false
		default:
			fmt.Println("Неверный выбор")
		}
	}
}

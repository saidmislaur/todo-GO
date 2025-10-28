package ui

import (
	"demo/app/internal/tasks"
	"fmt"
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
		case "1":
			tm.ShowTasks()
		case "2":
			tm.AddTask()
		case "3":
			tm.DeleteTask()
		case "4":
			tm.UpdateTask()
		case "5":
			return false
		default:
			fmt.Println("Неверный выбор")
		}
	}
}

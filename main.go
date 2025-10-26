package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type TaskManager struct {
	Tasks  []Task
	Reader *bufio.Reader
}

func main() {
	app := TaskManager{
		Reader: bufio.NewReader(os.Stdin),
	}

	app.loadTasks()

	app.Run()
}

func (tm *TaskManager) Run() {
	for {
		if !tm.showMenu() {
			fmt.Println("Программа завершена")
			break
		}
	}
}

func (tm *TaskManager) showMenu() bool {
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
		tm.showTasks()
	case "2":
		tm.addTask()
	case "3":
		tm.deleteTask()
	case "4":
		tm.updateTask()
	case "5":
		return false
	default:
		fmt.Println("Неверный выбор")
		return false
	}

	return true
}

func (tm *TaskManager) showTasks() {
	fmt.Println("\nСписок задач:")
	for _, task := range tm.Tasks {
		fmt.Printf("%d. %s\n", task.ID, task.Text)
	}
	fmt.Println()
}

func (tm *TaskManager) addTask() {
	fmt.Println("Введите новую задачу: ")
	text, _ := tm.Reader.ReadString('\n')
	text = strings.TrimSpace(text)

	newTask := Task{
		ID:   len(tm.Tasks) + 1,
		Text: text,
	}

	tm.Tasks = append(tm.Tasks, newTask)
	tm.saveTasks()
	fmt.Println("Задача добавлена ✅")
}

func (tm *TaskManager) deleteTask() {
	fmt.Println("Введите номер задачи для удаления: ")
	text, _ := tm.Reader.ReadString('\n')
	text = strings.TrimSpace(text)

	id, err := strconv.Atoi(text)
	if err != nil || id < 1 || id > len(tm.Tasks) {
		fmt.Println("некорректный номер задачи")
	}

	num := -1

	for i, task := range tm.Tasks {
		if task.ID == id {
			num = i
		}
	}

	if num == -1 {
		fmt.Println("Такой задачи не существует")
	}

	deleted := tm.Tasks[num]
	tm.Tasks = append(tm.Tasks[:num], tm.Tasks[num+1:]...)
	tm.saveTasks()

	fmt.Println("Удалена задача:", deleted.Text)
}

func (tm *TaskManager) updateTask() {
	fmt.Println("Введите номер задачи, которую хотите редактировать")
	text, _ := tm.Reader.ReadString('\n')

	text = strings.TrimSpace(text)
	num, err := strconv.Atoi(text)
	if err != nil || num < 1 || num > len(tm.Tasks) {
		fmt.Println("Введен некорректный номер задачи")
	}

	fmt.Println("Введите новый текст")
	newText, _ := tm.Reader.ReadString('\n')
	newText = strings.TrimSpace(newText)

	tm.Tasks[num-1].Text = newText
	tm.saveTasks()
	fmt.Println("Задача успешно обновлена ✅")
}

func (tm *TaskManager) saveTasks() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Ошибка при сохранении в файл", file)
	}

	json.NewEncoder(file).Encode(tm.Tasks)
}

func (tm *TaskManager) loadTasks() {
	file, err := os.Open("tasks.json")
	if err != nil {
		fmt.Println("Файл задач не найден, создаю новый")
		return
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&tm.Tasks)
}

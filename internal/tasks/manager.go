package tasks

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func NewManager() *TaskManager {
	return &TaskManager{
		Reader: bufio.NewReader(os.Stdin),
	}
}

func (tm *TaskManager) ShowTasks() {
	if len(tm.Tasks) == 0 {
		fmt.Println("Список задач пуст!")
	}
	fmt.Println("\nСписок задач:")
	for _, task := range tm.Tasks {
		fmt.Printf("%d. %s\n", task.ID, task.Text)
	}
	fmt.Println()
}

func (tm *TaskManager) AddTask() {
	fmt.Println("Введите новую задачу: ")
	text, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	text = strings.TrimSpace(text)

	newTask := Task{
		ID:   rand.Intn(100),
		Text: text,
	}

	tm.Tasks = append(tm.Tasks, newTask)
	tm.SaveTasks()
	fmt.Println("Задача добавлена ✅")
}

func (tm *TaskManager) DeleteTask() {
	fmt.Println("Введите номер для удаления: ")
	text, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	text = strings.TrimSpace(text)

	id, err := strconv.Atoi(text)
	if err != nil {
		log.Fatalf("%v", err)
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
	tm.SaveTasks()

	fmt.Println("Удалена задача:", deleted.Text)
}

func (tm *TaskManager) UpdateTask() {
	fmt.Println("Введите номер задачи, которую хотите редактировать")
	text, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	text = strings.TrimSpace(text)

	id, err := strconv.Atoi(text)
	if err != nil {
		log.Fatalf("%v", err)
	}

	num := -1

	for i, task := range tm.Tasks {
		if task.ID == id {
			num = i
		}
	}

	fmt.Println("Введите новый текст")
	newText, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	newText = strings.TrimSpace(newText)

	tm.Tasks[num].Text = newText
	tm.SaveTasks()
	fmt.Println("Задача успешно обновлена ✅")
}

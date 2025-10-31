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

func NewManager(filePath string) *TaskManager {
	return &TaskManager{
		Reader:   bufio.NewReader(os.Stdin),
		FilePath: filePath,
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

func (tm *TaskManager) AddTask() error {
	fmt.Println("Введите новую задачу: ")
	text, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("текст задачи не может быть пустым")
	}

	newTask := Task{
		ID:   rand.Intn(1000),
		Text: text,
	}

	tm.Tasks = append(tm.Tasks, newTask)
	if err := tm.SaveTasks(); err != nil {
		return fmt.Errorf("ошибка при добавлении задачи %w", err)
	}
	fmt.Println("Задача добавлена ✅")

	return nil
}

func (tm *TaskManager) DeleteTask() error {
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
		return fmt.Errorf("Такой задачи не существует", id)
	}

	deleted := tm.Tasks[num]
	tm.Tasks = append(tm.Tasks[:num], tm.Tasks[num+1:]...)
	if err := tm.SaveTasks(); err != nil {
		return fmt.Errorf("ошибка при сохранении после удаления: %w", err)
	}

	fmt.Println("Удалена задача:", deleted.Text)

	return err
}

func (tm *TaskManager) UpdateTask() error {
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
			break
		}
	}

	if num == -1 {
		fmt.Println("Такой задачи не существует")
	}

	fmt.Println("Введите новый текст")
	newText, err := tm.Reader.ReadString('\n')
	if err != nil {
		log.Fatalf("%v", err)
	}
	newText = strings.TrimSpace(newText)

	tm.Tasks[num].Text = newText
	if err := tm.SaveTasks(); err != nil {
		return fmt.Errorf("ошибка сохранения изменений: %w", err)
	}
	fmt.Println("Задача успешно обновлена ✅")

	return nil
}

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
		fmt.Printf("%d. %s – [%v]\n", task.ID, task.Text, task.Status)
	}
	fmt.Println()
}

func (tm *TaskManager) AddTask() error {
	fmt.Println("Введите новую задачу: ")
	text, err := tm.Reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%v", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("текст задачи не может быть пустым")
	}

	fmt.Println("статус задачи по умолчанию: in_process")
	status := "in_process"

	newTask := Task{
		ID:     rand.Intn(1000),
		Text:   text,
		Status: status,
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
		fmt.Println("Такой задачи не существует:", id)
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
		return fmt.Errorf("задача с ID %d не существует", id)

	}

	fmt.Println("1. Редактировать текст")
	fmt.Println("2. Редактировать статус")
	fmt.Println("3. Отмена")

	var variant string
	if _, err := fmt.Scan(&variant); err != nil {
		return fmt.Errorf("ошибка чтения варианта: %w", err)
	}

	switch variant {
	case "1":
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

	case "2":
		fmt.Println("Обновите статус задачи: done / in_process / not_completed")
		status, err := tm.Reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		status = strings.TrimSpace(status)
		if status == "" && status != "done" && status != "in_process" && status != "not_completed" {
			return fmt.Errorf("недопустимый статус: %q", status)
		}

		tm.Tasks[num].Status = status

		fmt.Printf("Статус задачи %v изменён на %v ✅", tm.Tasks[num].Text, status)
	case "3":
		break
	}
	return nil
}

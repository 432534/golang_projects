package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var tasks []Task

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func main() {
	fmt.Println("Welcome to the Task Tracker!")
	err := loadTask()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Show tasks")
		fmt.Println("2. Add a task")
		fmt.Println("3. Update task status")
		fmt.Println("4. Exit")
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			showTask()
		case "2":
			fmt.Print("Enter task description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)
			addTask(description)
		case "3":
			updatetaskstatus()
		case "4":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}
func loadTask() error {
	fileData, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return nil
		}
		return err
	}
	return json.Unmarshal(fileData, &tasks)
}
func saveTask() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644)
}
func addTask(description string) {
	now := time.Now().Unix()
	newTask := Task{
		ID:          getNextID(),
		Description: description,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks = append(tasks, newTask)
	err := saveTask()
	if err != nil {
		fmt.Println("Error saving tasks:", err)
	} else {
		fmt.Println("Task added successfully!")
	}
}
func showTask() {
	if len(tasks) == 0 {
		fmt.Println("There are no tasks available.")
		return
	}
	fmt.Println("\nYour Tasks:")
	for _, task := range tasks {
		fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreated At: %s\nUpdated At: %s\n\n",
			task.ID,
			task.Description,
			task.Status,
			time.Unix(task.CreatedAt, 0).Format(time.RFC822),
			time.Unix(task.UpdatedAt, 0).Format(time.RFC822),
		)
	}
}
func updatetaskstatus() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the ID of the task you want to update: ")
	var idx int
	_, err := fmt.Scanf("%d\n", &idx)
	if err != nil {
		fmt.Println("Enter a valid numeric ID.")
		return
	}

	found := false
	for i, task := range tasks {
		if task.ID == idx {
			fmt.Printf("The current status of the task is: %s\n", task.Status)
			fmt.Print("Enter the new status (pending/in-progress/completed): ")
			newstatus, _ := reader.ReadString('\n')
			newstatus = strings.TrimSpace(newstatus)

			if newstatus != "pending" && newstatus != "in-progress" && newstatus != "completed" {
				fmt.Println("Invalid status. Use: pending, in-progress, or completed.")
				return
			}

			tasks[i].Status = newstatus
			tasks[i].UpdatedAt = time.Now().Unix()

			if err := saveTask(); err != nil {
				fmt.Println("Error saving task:", err)
			} else {
				fmt.Println("Task status updated successfully.")
			}

			found = true
			break
		}
	}

	if !found {
		fmt.Println("Task with the given ID not found.")
	}
}

func getNextID() int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

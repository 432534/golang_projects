package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Expense struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

const dataFile = "expenses.json"

func loadExpense() ([]Expense, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return []Expense{}, nil
	}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		panic(err)
	}
	var expenses []Expense
	err = json.Unmarshal(data, &expenses)
	return expenses, err
}
func saveExpenses(expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}
func addExpense(description string, amount float64) error {
	if amount < 0 {
		fmt.Println("Please enter the valid value in the amount part")
	}

	expenses, err := loadExpense()
	if err != nil {
		panic(err)
	}

	newID := 1
	if len(expenses) > 0 {
		newID = expenses[len(expenses)-1].ID + 1
	}

	expense := Expense{
		ID:          newID,
		Date:        time.Now().Format("2006-01-02"),
		Description: description,
		Amount:      amount,
	}

	expenses = append(expenses, expense)

	err = saveExpenses(expenses)
	if err != nil {
		return err
	}

	fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)
	return nil
}

func main() {
	fmt.Println("We are going to develop CLi for expense-tracker")
	if len(os.Args) < 2 || os.Args[1] != "add" {
		println("Provide a valid argument")
	}
	addflag := flag.NewFlagSet("add", flag.ExitOnError)
	desc := addflag.String("description", "", "Expense description")
	amt := addflag.Float64("amount", 0, "Expense amount")

	addflag.Parse(os.Args[2:])
	if *desc == "" || *amt == 0 {
		fmt.Println("Both --description and --amount are required")
		return
	}
	err := addExpense(*desc, *amt)
	if err != nil {
		fmt.Println("Error adding expense:", err)
	}
}

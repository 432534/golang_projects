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
func updateExpense(description string, amount float64, Id int) error {
	expenses, err := loadExpense()
	if err != nil {
		panic(err)
	}
	update := false
	for i, exp := range expenses {
		if exp.ID == Id {
			if description != "" {
				expenses[i].Description = description
			}
			if amount > 0 {
				expenses[i].Amount = amount
			}
			update = true
			break
		}
	}
	if !update {
		fmt.Println("the Id that you are searching is not available")
	}
	err = saveExpenses(expenses)
	if err != nil {
		return err
	}

	fmt.Println("Expense updated successfully")
	return nil
}
func deleteExpense(Id int) error {
	expenses, err := loadExpense()
	if err != nil {
		panic(err)
	}
	found := false
	var newExpense = []Expense{}
	for _, exp := range expenses {
		if exp.ID == Id {
			found = true
			continue
		}
		newExpense = append(newExpense, exp)
	}

	if !found {
		return fmt.Errorf("no expense found with ID: %d", Id)
	}
	err = saveExpenses(newExpense)
	if err != nil {
		fmt.Println("Error in saving the deleted expenses")
	}
	return nil
}
func viewExpense() error {
	expenses, err := loadExpense()
	if err != nil {
		panic(err)
	}
	if len(expenses) == 0 {
		fmt.Println("there are no expenses to view")
	}
	for _, exp := range expenses {
		fmt.Printf("Id: %d | description: %s | amount: %.2f | date: %s\n", exp.ID, exp.Description, exp.Amount, exp.Date)
	}
	return nil
}
func viewSummayofExpense() error {
	expenses, err := loadExpense()
	if err != nil {
		panic(err)
	}
	if len(expenses) == 0 {
		fmt.Println("you have'nt spent any money yet")
	}
	total := 0.0
	for _, exp := range expenses {
		total += exp.Amount
	}
	fmt.Printf("The total expenses that you have done is:  ₹%.2f\n", total)
	return nil
}
func main() {
	fmt.Println("We are going to develop CLi for expense-tracker")
	if os.Args[1] == "add" {
		if len(os.Args) < 2 {
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
	if os.Args[1] == "update" {
		updateflag := flag.NewFlagSet("update", flag.ExitOnError)
		Id := updateflag.Int("id", 0, "ID to update")
		newdesc := updateflag.String("description", " ", "New description")
		newamt := updateflag.Float64("amount", 0, "New amount")
		updateflag.Parse(os.Args[2:])
		if *newdesc == "" || *newamt == 0 || *Id < 0 {
			fmt.Println("Both --description and --amount are required")
			return
		}
		err := updateExpense(*newdesc, *newamt, *Id)
		if err != nil {
			fmt.Println("Error in updating the expense", err)
		}
	}
	if os.Args[1] == "delete" {
		deleteflag := flag.NewFlagSet("delete", flag.ExitOnError)
		Id := deleteflag.Int("id", 0, "id to delete")
		if *Id < 0 {
			fmt.Println("Please provide the valid id for deletion")
		}
		deleteflag.Parse(os.Args[2:])
		err := deleteExpense(*Id)
		if err != nil {
			fmt.Println("Error in deleting the expense", err)
		}
	}
	if os.Args[1] == "view" {
		err := viewExpense()
		if err != nil {
			fmt.Println("there is an error in viewing the expenses", err)
		}
	}
	if os.Args[1] == "summary" {
		err := viewSummayofExpense()
		if err != nil {
			fmt.Println("there is an error in viewing the summay of the expenses", err)
		}
	}
}

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("Please select the difficulty level:")
	fmt.Println("1. Easy (10 chances)")
	fmt.Println("2. Medium (5 chances)")
	fmt.Println("3. Hard (3 chances)")

	var difficulty int
	fmt.Print("Enter difficulty (1/2/3): ")
	_, err := fmt.Scanln(&difficulty)
	if err != nil {
		fmt.Println("Entered number is invalid. Please enter a valid number.")
		return
	}

	chances := 5
	switch difficulty {
	case 1:
		chances = 10
	case 2:
		chances = 5
	case 3:
		chances = 3
	default:
		fmt.Println("Invalid difficulty choice. Please run the game again.")
		return
	}

	randomNumber := rand.Intn(100) + 1
	for i := 1; i <= chances; i++ {
		var guess int
		fmt.Printf("Guess #%d: ", i)
		_, err := fmt.Scanln(&guess)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			i--
			continue
		}

		if guess == randomNumber {
			fmt.Println("🎉 Congratulations! You guessed the correct number!")
			return
		} else if guess < randomNumber {
			fmt.Println("Too low. Try a higher number.")
		} else {
			fmt.Println("Too high. Try a lower number.")
		}
	}

	fmt.Printf("😢 Sorry, you ran out of chances. The number was %d.\n", randomNumber)
}

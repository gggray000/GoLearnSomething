package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const accountBalanceFile = "balance.txt"

func main() {
	var accountBalance, err = getBalanceFromFile(accountBalanceFile)

	if err != nil {
		panic(err)
	}

	fmt.Println("Welcome to Go Bank!")

	for {
		fmt.Println("What do you want to do?")
		fmt.Println("1. Check balance")
		fmt.Println("2. Deposit money")
		fmt.Println("3. Withdraw money")
		fmt.Println("4. Exit")

		var choice int
		fmt.Print("Your choice: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			fmt.Println("Your balance is:", accountBalance)

		case 2:
			fmt.Print("Your deposit: ")
			var depositAmount float64
			fmt.Scan(&depositAmount)

			if depositAmount <= 0 {
				fmt.Println("Invalid amount.")
				continue
			}

			accountBalance += depositAmount
			fmt.Println("Your updated balance is:", accountBalance)
			writeBalanceToFile(accountBalance)

		case 3:
			fmt.Print("Your withdrawal: ")
			var withdrawAmount float64
			fmt.Scan(&withdrawAmount)

			if withdrawAmount <= 0 || withdrawAmount > accountBalance {
				fmt.Println("Invalid amount.")
				continue
			}

			accountBalance -= withdrawAmount
			fmt.Println("Your updated balance is:", accountBalance)
			writeBalanceToFile(accountBalance)

		default:
			fmt.Println("Goodbye!")
			fmt.Println("Thanks for choosing Go Bank!")
			return
		}
	}
}

func writeBalanceToFile(balance float64) {
	balanceText := fmt.Sprint(balance)
	os.WriteFile(accountBalanceFile, []byte(balanceText), 0644)
}

func getBalanceFromFile(fileName string) (float64, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return 0, errors.New("Failed to find file ")
	}

	balanceText := string(data)
	balance, err := strconv.ParseFloat(balanceText, 64)

	if err != nil {
		return 0, errors.New("Failed to parse float ")
	}
	return balance, nil
}

package main

import (
	"Go-Basics/bank/util"
	"fmt"
	"github.com/Pallinder/go-randomdata"
)

const accountBalanceFile = "balance.txt"

func main() {
	var accountBalance, err = util.GetFloatValueFromFile(accountBalanceFile)

	if err != nil {
		panic(err)
	}

	fmt.Println("Welcome to Go Bank!")
	fmt.Println("Reach us 24/7 via", randomdata.PhoneNumber())

	for {
		presentOptions()

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
			util.WriteFloatValueToFile(accountBalance, accountBalanceFile)

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
			util.WriteFloatValueToFile(accountBalance, accountBalanceFile)

		default:
			fmt.Println("Goodbye!")
			fmt.Println("Thanks for choosing Go Bank!")
			return
		}
	}
}

package main

import "fmt"

func main() {
	var accountBalance float64 = 10000

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

		/*	if choice == 1 {
				fmt.Println("Your balance is:", accountBalance)
			} else if choice == 2 {
				fmt.Print("Your deposit: ")
				var depositAmount float64
				fmt.Scan(&depositAmount)

				if depositAmount <= 0 {
					fmt.Println("Invalid amount.")
					continue
				}

				accountBalance += depositAmount
				fmt.Println("Your updated balance is:", accountBalance)
			} else if choice == 3 {
				fmt.Print("Your withdrawal: ")
				var withdrawAmount float64
				fmt.Scan(&withdrawAmount)

				if withdrawAmount <= 0 || withdrawAmount > accountBalance {
					fmt.Println("Invalid amount.")
					continue
				}

				accountBalance -= withdrawAmount
				fmt.Println("Your updated balance is:", accountBalance)
			} else {
				fmt.Println("Goodbye!")
				break
			}
		}*/

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

		default:
			fmt.Println("Goodbye!")
			fmt.Println("Thanks for choosing Go Bank!")
			return
		}
	}
}

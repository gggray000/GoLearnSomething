package main

import "fmt"

func main() {

	revenue := getUserInput("Revenue: ")
	expenses := getUserInput("Expenses: ")
	tax_rate := getUserInput("Tax Rate: ")

	ebt, profit, ratio := calculate(revenue, expenses, tax_rate)
	fmt.Printf("EBT: %.2f\n", ebt)
	fmt.Printf("Profit: %.2f\n", profit)
	fmt.Printf("Ratio: %.2f\n", ratio)
}

func getUserInput(infoText string) float64 {
	var userInput float64
	fmt.Print(infoText)
	fmt.Scan(&userInput)
	return userInput
}

func calculate(revenue, expenses, tax_rate float64) (float64, float64, float64) {
	ebt := revenue - expenses
	profit := ebt * (1 - tax_rate/100)
	ratio := ebt / profit
	return ebt, profit, ratio
}

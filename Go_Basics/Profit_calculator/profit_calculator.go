package main

import (
	"errors"
	"fmt"
	"os"
)

const resultFile = "result.txt"

func main() {

	revenue, err := getUserInput("Revenue: ")
	catchErr(err)
	expenses, err2 := getUserInput("Expenses: ")
	catchErr(err2)
	tax_rate, err3 := getUserInput("Tax Rate: ")
	catchErr(err3)

	ebt, profit, ratio := calculate(revenue, expenses, tax_rate)
	fmt.Printf("EBT: %.2f\n", ebt)
	fmt.Printf("Profit: %.2f\n", profit)
	fmt.Printf("Ratio: %.2f\n", ratio)
	writeToFile(ebt, profit, ratio)
}

func getUserInput(infoText string) (float64, error) {
	var userInput float64
	fmt.Print(infoText)
	fmt.Scan(&userInput)

	if userInput <= 0 {
		return 0, errors.New("Invalid Input ")
	}

	return userInput, nil
}

func catchErr(err error) {
	if err != nil {
		panic(err)
	}
}

func calculate(revenue, expenses, tax_rate float64) (float64, float64, float64) {
	ebt := revenue - expenses
	profit := ebt * (1 - tax_rate/100)
	ratio := ebt / profit
	return ebt, profit, ratio
}

func writeToFile(ebt, profit, ratio float64) {
	results := fmt.Sprintf("EBT: %.1f\nProfit: %.1f\nRation: %.2f\n",
		ebt, profit, ratio)
	os.WriteFile(resultFile, []byte(results), 0644)
}

package main

import "fmt"

func main() {
	var revenue float64
	var expenses float64
	var tax_rate float64

	fmt.Print("Revenue: ")
	fmt.Scan(&revenue)

	fmt.Print("Expenses: ")
	fmt.Scan(&expenses)

	fmt.Print("Tax Rate: ")
	fmt.Scan(&tax_rate)

	ebt := revenue - expenses
	profit := ebt * (1 - tax_rate/100)
	ratio := ebt / profit
	fmt.Print("EBT: ", ebt, " Profit: ", profit, " Ratio: ", ratio)
}

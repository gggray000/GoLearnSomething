package main

import (
	"fmt"
	"math"
)

const inflationRate = 2.0

func main() {
	var investmentAmount float64
	var expectedReturnRate float64
	var years float64

	fmt.Print("Investment Value: ")
	fmt.Scan(&investmentAmount)

	fmt.Print("Return Rate: ")
	fmt.Scan(&expectedReturnRate)

	fmt.Print("Years: ")
	fmt.Scan(&years)

	futureValue, futureRealValue := calculateFutureValues(investmentAmount, expectedReturnRate, years)

	//fmt.Println("Future Value: ", futureValue)
	fmt.Printf("Future Value: %.2f\nFuture Value(adjusted for inflation): %.2f",
		futureValue, futureRealValue)
}

func calculateFutureValues(investmentAmount, expectedReturnRate, years float64) (float64, float64) {
	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/100, years)
	futureRealValue := futureValue / math.Pow(1+inflationRate/100, years)
	return futureValue, futureRealValue
}

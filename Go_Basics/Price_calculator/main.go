package main

import (
	"Go-Basics/price-calculator/prices"
	"Go-Basics/price-calculator/util"
	"fmt"
)

func main() {

	taxRates := []float64{0, 0.07, 0.1, 0.15}

	for _, taxRate := range taxRates {
		fm := util.NewFileManager("prices.txt", fmt.Sprintf("result_%.0f.json", taxRate*100))
		//cmdm := util.NewCmdManager()
		priceJob := prices.New(*fm, taxRate)
		err := priceJob.Process()
		if err != nil {
			fmt.Println("Could not process job.")
			panic(err)
		}
	}

}

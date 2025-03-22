package main

import (
	"Go-Basics/price-calculator/prices"
	"Go-Basics/price-calculator/util"
	"fmt"
)

func main() {

	taxRates := []float64{0, 0.07, 0.1, 0.15}
	doneChans := make([]chan bool, len(taxRates))
	errorChans := make([]chan error, len(taxRates))

	for index, taxRate := range taxRates {
		doneChans[index] = make(chan bool)
		errorChans[index] = make(chan error)
		fm := util.NewFileManager("prices.txt", fmt.Sprintf("result_%.0f.json", taxRate*100))
		//cmdm := util.NewCmdManager()
		priceJob := prices.New(*fm, taxRate)
		go priceJob.Process(doneChans[index], errorChans[index])
		//if err != nil {
		//	fmt.Println("Could not process job.")
		//	panic(err)
		//}
	}

	for index, _ := range taxRates {
		select {
		case err := <-errorChans[index]:
			if err != nil {
				fmt.Println(err)
			}
		case <-doneChans[index]:
			fmt.Println("Done!")
		}

	}

}

package util

import "fmt"

type CmdManager struct {
}

func NewCmdManager() CmdManager {
	return CmdManager{}
}

func (cmd CmdManager) ReadLines() ([]string, error) {
	fmt.Println("Enter prices, confirm with ENTER")
	var prices []string
	for {
		var price string
		fmt.Println("Price: ")
		fmt.Scan(&price)

		if price == "0" {
			break
		}

		prices = append(prices, price)
	}
	return prices, nil
}

func (cmd CmdManager) WriteResult(data interface{}) error {
	fmt.Println(data)
	return nil
}

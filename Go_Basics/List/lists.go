package main

import "fmt"

type Product struct {
	id    string
	title string
	price float64
}

func main() {
	hobbies := [3]string{"eat", "drink", "play"}
	fmt.Println(hobbies[1:])

	mainHobbies := hobbies[:2]
	fmt.Println(mainHobbies)

	fmt.Println(cap(mainHobbies))
	secondaryHobbies := mainHobbies[1:3]
	fmt.Println(secondaryHobbies)

	courseGoals := []string{"Learb Go", "Learn Gin"}
	courseGoals[1] = "Learn Go 2"
	courseGoals = append(courseGoals, "Learn Go Zero and Gin")
	fmt.Println(courseGoals)

	products := []Product{
		{"1", "iPhone", 499},
		{"2", "Nike", 99},
	}
	products = append(products, Product{"3", "Coffee", 2.99})
	fmt.Println(products)

	prices := []float64{14.99, 10.99}
	discountPrices := []float64{9.99, 12.99, 18.99}
	prices = append(prices, discountPrices...)
	fmt.Println(prices)

}

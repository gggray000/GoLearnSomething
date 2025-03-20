package main

import "fmt"

func main() {
	//numbers := []int{1, 2, 3, 4}
	//doubled := transformNumbers(&numbers, double)
	//fmt.Println(doubled)
	//fmt.Println(factorial(6))
	numbers := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(sum(1, 2, 3, 4, 5))
	fmt.Println(sum(numbers...))
}

//func transformNumbers(numbers *[]int, transform func(int) int) []i nt {
//	dNumbers := []int{}
//	for _, value := range *numbers {
//		dNumbers = append(dNumbers, transform(value))
//	}
//	return dNumbers
//}
//
//func double(number int) int {
//	return number * 2
//}
//
//func triple(number int) int {
//	return number * 3
//}
//
//func getTransformerFunction() func(int) int {
//	return double
//}

//func factorial(number int) int {
//	if number == 1 {
//		return number
//	}
//	number *= factorial(number - 1)
//	return number
//}

func sum(numbers ...int) int {
	sum := 0
	for _, val := range numbers {
		sum += val
	}
	return sum
}

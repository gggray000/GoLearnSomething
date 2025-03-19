package main

import "fmt"

func main() {
	age := 32

	agePointer := &age

	fmt.Println("Age:", *agePointer)

	editAgeToAdultYears(agePointer)
	fmt.Println(age)
}

func editAgeToAdultYears(agePointer *int) {
	*agePointer = *agePointer - 18
}

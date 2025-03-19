package main

import (
	"Go-Basics/structs/user"
	"fmt"
)

func main() {
	firstName := getUserData("Please enter your first name: ")
	lastName := getUserData("Please enter your last name: ")
	birthdate := getUserData("Please enter your birthdate (MM/DD/YYYY): ")

	var appUser *user.User
	appUser, err := user.New(firstName, lastName, birthdate)

	if err != nil {
		panic(err)
	}

	appUser.OutputUserDetail()
	appUser.ClearUserName()
	appUser.OutputUserDetail()

	admin, _ := user.NewAdmin("admin@test.com", "test123")
	admin.OutputUserDetail()

}

func getUserData(promptText string) string {
	fmt.Print(promptText)
	var value string
	fmt.Scanln(&value)
	return value
}

package user

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	firstName string
	lastName  string
	birthdate string
	createdAt time.Time
}

type Admin struct {
	email    string
	password string
	User
}

func NewAdmin(email, password string) (*Admin, error) {
	if email == "" || password == "" {
		return nil, errors.New("No null value ")
	}

	return &Admin{
		email:    email,
		password: password,
		User: User{
			firstName: "Admin",
			lastName:  "Admin",
			birthdate: "---",
			createdAt: time.Now(),
		},
	}, nil
}

func New(firstName, lastName, birthdate string) (*User, error) {

	if firstName == "" || lastName == "" || birthdate == "" {
		return nil, errors.New("No null value ")
	}

	return &User{
		firstName: firstName,
		lastName:  lastName,
		birthdate: birthdate,
		createdAt: time.Now(),
	}, nil
}

func (u *User) OutputUserDetail() {
	fmt.Println(u.firstName, u.lastName, u.birthdate)
}

func (u *User) ClearUserName() {
	u.firstName = ""
	u.lastName = ""
}

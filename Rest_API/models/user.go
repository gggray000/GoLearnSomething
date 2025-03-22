package models

import (
	"errors"
	"rest-api/database"
	"rest-api/util"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"
	statement, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := statement.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return err
	}
	isPasswordValid := util.CheckPasswordHash(u.Password, retrievedPassword)
	if !isPasswordValid {
		return errors.New("Invalid Credentials! ")
	}
	return nil
}

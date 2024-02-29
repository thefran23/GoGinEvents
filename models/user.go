package models

import (
	"errors"
	"fmt"

	"example.com/db"
	"example.com/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hash)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()

	return err

}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users where email = ?"
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}
	fmt.Println("u.Password")
	fmt.Println(u.Password)

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}
	return nil
}

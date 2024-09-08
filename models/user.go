package models

import (
	"database/sql"
	"encoding/json"
	"errors"

	"example.com/db"
	"example.com/utils"
)

type User struct {
	ID       int64
	Login    string   `binding:"required"`
	Password string   `binding:"required"`
	Roles    []string `json:"roles"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(login, password) VALUES ($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = stmt.QueryRow(u.Login, hashedPassword).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password, roles FROM users WHERE login = $1"
	row := db.DB.QueryRow(query, u.Login)

	var retrievedPassword string
	var roles []byte
	err := row.Scan(&u.ID, &retrievedPassword, &roles)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Credentials invalid")
		}
		return err
	}

	err = json.Unmarshal(roles, &u.Roles)
	if err != nil {
		return errors.New("Failed to parse roles")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}

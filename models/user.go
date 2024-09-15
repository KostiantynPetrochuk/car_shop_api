package models

import (
	"database/sql"
	"errors"

	"example.com/db"
	"example.com/utils"
	"github.com/lib/pq"
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
	var roles []string
	err := row.Scan(&u.ID, &retrievedPassword, pq.Array(&roles))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("credentials invalid")
		}
		return err
	}

	u.Roles = roles

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}

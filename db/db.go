package db

import (
	"database/sql"
	"fmt"

	"example.com/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := config.GetEnv("DATABASE_URL")
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Could not connected to db.")
	}
	// DB.SetMaxOpenConns(10)
	// DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		roles TEXT[]
	)
	`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created users table.")
	}

	createBrandsTable := `
	CREATE TABLE IF NOT EXISTS brands (
		id SERIAL PRIMARY KEY,
		brand_name TEXT NOT NULL UNIQUE,
		file_name TEXT NOT NULL
	)
	`

	_, err = DB.Exec(createBrandsTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created brands table.")
	}

	createModelsTable := `
	CREATE TABLE IF NOT EXISTS models (
		id SERIAL PRIMARY KEY,
		model_name TEXT NOT NULL,
		brand_id INTEGER NOT NULL,
		FOREIGN KEY (brand_id) REFERENCES brands(id)
	)
	`

	_, err = DB.Exec(createModelsTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created models table.")
	}
}

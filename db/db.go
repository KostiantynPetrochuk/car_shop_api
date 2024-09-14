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

	createFeaturesTable := `
	CREATE TABLE IF NOT EXISTS features (
		id SERIAL PRIMARY KEY,
		feature_name TEXT NOT NULL,
		category TEXT NOT NULL
	)
	`

	_, err = DB.Exec(createFeaturesTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created features table.")
	}

	createCarsTable := `
	CREATE TABLE IF NOT EXISTS cars (
    	id SERIAL PRIMARY KEY,
    	vin TEXT NOT NULL UNIQUE,
    	brand_id INTEGER NOT NULL,
    	model_id INTEGER NOT NULL,
    	body_type TEXT NOT NULL,
    	fuel_type TEXT NOT NULL,
    	year INTEGER NOT NULL,
    	transmission TEXT NOT NULL,
    	drive_type TEXT NOT NULL,
    	condition TEXT NOT NULL,
    	engine_size FLOAT NOT NULL,
    	door_count INTEGER NOT NULL,
    	cylinder_count INTEGER NOT NULL,
    	color TEXT NOT NULL,
    	mileage INTEGER NOT NULL,
    	image_names JSONB,
    	FOREIGN KEY (brand_id) REFERENCES brands(id),
    	FOREIGN KEY (model_id) REFERENCES models(id)
	)
	`

	_, err = DB.Exec(createCarsTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created cars table.")
	}

	createCarFeaturesTable := `
	CREATE TABLE IF NOT EXISTS car_features (
		id SERIAL PRIMARY KEY,
		car_id INTEGER NOT NULL,
		feature_id INTEGER NOT NULL,
		FOREIGN KEY (car_id) REFERENCES cars(id),
		FOREIGN KEY (feature_id) REFERENCES features(id)
	)
	`

	_, err = DB.Exec(createCarFeaturesTable)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Could not created car_features table.")
	}

}

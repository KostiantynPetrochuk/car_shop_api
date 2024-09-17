package models

import (
	"encoding/json"
	"log"

	"example.com/db"
)

type Car struct {
	ID           int64
	VIN          string  `binding:"required"`
	BrandId      int     `binding:"required"`
	ModelId      int     `binding:"required"`
	BodyType     string  `binding:"required"`
	Mileage      int     `binding:"required"`
	FuelType     string  `binding:"required"`
	Year         int     `binding:"required"`
	Transmission string  `binding:"required"`
	DriveType    string  `binding:"required"`
	Condition    string  `binding:"required"`
	EngineSize   float64 `binding:"required"`
	DoorCount    int     `binding:"required"`
	Price        int     `binding:"required"`
	Color        string  `binding:"required"`
	ImageNames   json.RawMessage
	BrandName    string `json:"BrandName"`
	ModelName    string `json:"ModelName"`
	CreatedAt    string `json:"CreatedAt"`
	Features     []Feature
}

func (c *Car) Save() error {
	query := `INSERT INTO cars(vin, brand_id, model_id, body_type, mileage, fuel_type, year, transmission, drive_type, condition, engine_size, door_count, price, color, image_names) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(c.VIN, c.BrandId, c.ModelId, c.BodyType, c.Mileage, c.FuelType, c.Year, c.Transmission, c.DriveType, c.Condition, c.EngineSize, c.DoorCount, c.Price, c.Color, c.ImageNames).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetCars() ([]Car, error) {
	var cars []Car
	query := `
		SELECT 
			cars.id, cars.vin, cars.brand_id, cars.model_id, cars.body_type, cars.mileage, cars.fuel_type, cars.year, 
			cars.transmission, cars.drive_type, cars.condition, cars.engine_size, cars.door_count, cars.price, 
			cars.color, cars.image_names, cars.created_at, brands.brand_name AS brand_name, models.model_name AS model_name
		FROM 
			cars
		JOIN 
			brands ON cars.brand_id = brands.id
		JOIN 
			models ON cars.model_id = models.id`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.VIN, &car.BrandId, &car.ModelId, &car.BodyType, &car.Mileage, &car.FuelType, &car.Year, &car.Transmission, &car.DriveType, &car.Condition, &car.EngineSize, &car.DoorCount, &car.Price, &car.Color, &car.ImageNames, &car.CreatedAt, &car.BrandName, &car.ModelName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return cars, nil
}

func GetCarById(id int64) (Car, error) {
	var car Car
	query := `
		SELECT 
			cars.id, cars.vin, cars.brand_id, cars.model_id, cars.body_type, cars.mileage, cars.fuel_type, cars.year, 
			cars.transmission, cars.drive_type, cars.condition, cars.engine_size, cars.door_count, cars.price, 
			cars.color, cars.image_names, cars.created_at, brands.brand_name AS brand_name, models.model_name AS model_name
		FROM 
			cars
		JOIN 
			brands ON cars.brand_id = brands.id
		JOIN 
			models ON cars.model_id = models.id
		WHERE 
			cars.id = $1`

	err := db.DB.QueryRow(query, id).Scan(&car.ID, &car.VIN, &car.BrandId, &car.ModelId, &car.BodyType, &car.Mileage, &car.FuelType, &car.Year, &car.Transmission, &car.DriveType, &car.Condition, &car.EngineSize, &car.DoorCount, &car.Price, &car.Color, &car.ImageNames, &car.CreatedAt, &car.BrandName, &car.ModelName)
	if err != nil {
		log.Println("Error scanning row:", err)
		return car, err
	}

	featuresQuery := `
		SELECT 
			features.id, features.feature_name, features.category
		FROM 
			features
		JOIN 
			car_features ON features.id = car_features.feature_id
		WHERE 
			car_features.car_id = $1`

	rows, err := db.DB.Query(featuresQuery, id)
	if err != nil {
		log.Println("Error executing features query:", err)
		return car, err
	}
	defer rows.Close()

	var features []Feature
	for rows.Next() {
		var feature Feature
		if err := rows.Scan(&feature.ID, &feature.FeatureName, &feature.Category); err != nil {
			log.Println("Error scanning feature:", err)
			return car, err
		}
		features = append(features, feature)
	}

	car.Features = features

	return car, nil
}

package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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

type CarFilter struct {
	Offset       int
	Limit        int
	Condition    string
	Brand        string
	Model        string
	BodyType     string
	MileageFrom  int
	MileageTo    int
	FuelType     string
	Transmission string
	DriveType    string
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

func GetCars(filter CarFilter) ([]Car, int, error) {
	var cars []Car
	var total int
	argIndex := 1

	countQuery := `SELECT COUNT(*) FROM cars`
	var countConditions []string

	if filter.Condition != "" {
		countConditions = append(countConditions, `condition IN (`+formatCondition(filter.Condition)+`)`)
	}

	var brandId int
	var err error
	if filter.Brand != "" {
		brandId, err = strconv.Atoi(filter.Brand)
		if err == nil {
			countConditions = append(countConditions, `brand_id = $`+strconv.Itoa(argIndex))
			argIndex++
		}
	}

	var modelIds []int
	if filter.Model != "" {
		modelParts := strings.Split(filter.Model, ",")
		for _, part := range modelParts {
			modelId, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil {
				modelIds = append(modelIds, modelId)
			}
		}
		if len(modelIds) > 0 {
			modelPlaceholders := make([]string, len(modelIds))
			for i := range modelIds {
				modelPlaceholders[i] = "$" + strconv.Itoa(argIndex)
				argIndex++
			}
			countConditions = append(countConditions, `model_id IN (`+strings.Join(modelPlaceholders, ",")+`)`)
		}
	}

	var bodyTypes []string
	if filter.BodyType != "" {
		bodyTypeParts := strings.Split(filter.BodyType, ",")
		for _, part := range bodyTypeParts {
			bodyTypes = append(bodyTypes, "'"+strings.TrimSpace(part)+"'")
		}
		if len(bodyTypes) > 0 {
			countConditions = append(countConditions, `body_type IN (`+strings.Join(bodyTypes, ",")+`)`)
		}
	}

	var fuelTypes []string
	if filter.FuelType != "" {
		fuelTypeParts := strings.Split(filter.FuelType, ",")
		for _, part := range fuelTypeParts {
			fuelTypes = append(fuelTypes, "'"+strings.TrimSpace(part)+"'")
		}
		if len(fuelTypes) > 0 {
			countConditions = append(countConditions, `fuel_type IN (`+strings.Join(fuelTypes, ",")+`)`)
		}
	}

	var transmissions []string
	if filter.Transmission != "" {
		transmissionParts := strings.Split(filter.Transmission, ",")
		for _, part := range transmissionParts {
			transmissions = append(transmissions, "'"+strings.TrimSpace(part)+"'")
		}
		if len(transmissions) > 0 {
			countConditions = append(countConditions, `transmission IN (`+strings.Join(transmissions, ",")+`)`)
		}
	}

	var driveTypes []string
	if filter.DriveType != "" {
		driveTypeParts := strings.Split(filter.DriveType, ",")
		for _, part := range driveTypeParts {
			driveTypes = append(driveTypes, "'"+strings.TrimSpace(part)+"'")
		}
		if len(driveTypes) > 0 {
			countConditions = append(countConditions, `drive_type IN (`+strings.Join(driveTypes, ",")+`)`)
		}
	}

	if filter.MileageFrom > 0 {
		countConditions = append(countConditions, `mileage >= $`+strconv.Itoa(argIndex))
		argIndex++
	}
	if filter.MileageTo > 0 {
		countConditions = append(countConditions, `mileage <= $`+strconv.Itoa(argIndex))
		argIndex++
	}

	if len(countConditions) > 0 {
		countQuery += ` WHERE ` + strings.Join(countConditions, ` AND `)
	}

	var countArgs []interface{}
	if filter.Brand != "" {
		countArgs = append(countArgs, brandId)
	}
	for _, modelId := range modelIds {
		countArgs = append(countArgs, modelId)
	}
	if filter.MileageFrom > 0 {
		countArgs = append(countArgs, filter.MileageFrom)
	}
	if filter.MileageTo > 0 {
		countArgs = append(countArgs, filter.MileageTo)
	}

	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Println("Error executing count query:", err)
		return nil, 0, err
	}

	argIndex = 3
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

	var conditions []string
	if filter.Condition != "" {
		conditions = append(conditions, `cars.condition IN (`+formatCondition(filter.Condition)+`)`)
	}
	if filter.Brand != "" {
		conditions = append(conditions, `cars.brand_id = $`+strconv.Itoa(argIndex))
		argIndex++
	}
	if len(modelIds) > 0 {
		modelPlaceholders := make([]string, len(modelIds))
		for i := range modelIds {
			modelPlaceholders[i] = "$" + strconv.Itoa(argIndex)
			argIndex++
		}
		conditions = append(conditions, `cars.model_id IN (`+strings.Join(modelPlaceholders, ",")+`)`)
	}
	if len(bodyTypes) > 0 {
		conditions = append(conditions, `cars.body_type IN (`+strings.Join(bodyTypes, ",")+`)`)
	}
	if len(fuelTypes) > 0 {
		conditions = append(conditions, `cars.fuel_type IN (`+strings.Join(fuelTypes, ",")+`)`)
	}
	if len(transmissions) > 0 {
		conditions = append(conditions, `cars.transmission IN (`+strings.Join(transmissions, ",")+`)`)
	}
	if len(driveTypes) > 0 {
		conditions = append(conditions, `cars.drive_type IN (`+strings.Join(driveTypes, ",")+`)`)
	}
	if filter.MileageFrom > 0 {
		conditions = append(conditions, `cars.mileage >= $`+strconv.Itoa(argIndex))
		argIndex++
	}
	if filter.MileageTo > 0 {
		conditions = append(conditions, `cars.mileage <= $`+strconv.Itoa(argIndex))
		argIndex++
	}

	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, ` AND `)
	}
	query += ` ORDER BY cars.created_at DESC LIMIT $2 OFFSET $1`

	var args []interface{}
	args = append(args, filter.Offset, filter.Limit)
	if filter.Brand != "" {
		args = append(args, brandId)
	}
	for _, modelId := range modelIds {
		args = append(args, modelId)
	}
	if filter.MileageFrom > 0 {
		args = append(args, filter.MileageFrom)
	}
	if filter.MileageTo > 0 {
		args = append(args, filter.MileageTo)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.VIN, &car.BrandId, &car.ModelId, &car.BodyType, &car.Mileage, &car.FuelType, &car.Year, &car.Transmission, &car.DriveType, &car.Condition, &car.EngineSize, &car.DoorCount, &car.Price, &car.Color, &car.ImageNames, &car.CreatedAt, &car.BrandName, &car.ModelName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, 0, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return cars, total, nil
}

func formatCondition(condition string) string {
	conditions := strings.Split(condition, ",")
	for i := range conditions {
		conditions[i] = "'" + strings.TrimSpace(conditions[i]) + "'"
	}
	return strings.Join(conditions, ",")
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

func GetFeaturedCars() ([]Car, []Car, error) {
	var intactCars []Car

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
			cars.condition = 'intact'
		ORDER BY cars.created_at DESC
		LIMIT 5`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		return nil, nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.VIN, &car.BrandId, &car.ModelId, &car.BodyType, &car.Mileage, &car.FuelType, &car.Year, &car.Transmission, &car.DriveType, &car.Condition, &car.EngineSize, &car.DoorCount, &car.Price, &car.Color, &car.ImageNames, &car.CreatedAt, &car.BrandName, &car.ModelName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, nil, err
		}
		intactCars = append(intactCars, car)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// damaged cars
	var damagedCars []Car

	secondQuery := `
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
			cars.condition = 'damaged'
		ORDER BY cars.created_at DESC
		LIMIT 5`

	second_stmt, err := db.DB.Prepare(secondQuery)
	if err != nil {
		log.Println("Error preparing query:", err)
		return nil, nil, err
	}
	defer stmt.Close()

	secondRows, err := second_stmt.Query()
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, nil, err
	}
	defer secondRows.Close()

	for secondRows.Next() {
		var car Car
		err := secondRows.Scan(&car.ID, &car.VIN, &car.BrandId, &car.ModelId, &car.BodyType, &car.Mileage, &car.FuelType, &car.Year, &car.Transmission, &car.DriveType, &car.Condition, &car.EngineSize, &car.DoorCount, &car.Price, &car.Color, &car.ImageNames, &car.CreatedAt, &car.BrandName, &car.ModelName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, nil, err
		}
		damagedCars = append(damagedCars, car)
	}

	if err = secondRows.Err(); err != nil {
		log.Fatal(err)
	}

	return intactCars, damagedCars, nil
}

func GetCarsByBrand(brandId int) ([]Car, error) {
	var cars []Car
	fmt.Println("Brand ID:", brandId)

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

		ORDER BY cars.created_at DESC
		LIMIT 10`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
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

	fmt.Println(cars)

	return cars, nil
}

package models

import (
	"database/sql"
	"fmt"
	"sort"

	"example.com/db"
)

type Brand struct {
	ID        int64
	BrandName string `binding:"required"`
	FileName  string `binding:"required"`
	Models    []Model
}

func (b *Brand) Save() (Brand, error) {
	query := "INSERT INTO brands(brand_name, file_name) VALUES ($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Brand{}, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(b.BrandName, b.FileName).Scan(&b.ID)
	if err != nil {
		return Brand{}, err
	}

	return *b, nil
}

func GetBrands() ([]Brand, error) {
	query := `
		SELECT b.id, b.brand_name, b.file_name, m.id, m.model_name, m.brand_id
		FROM brands b
		LEFT JOIN models m ON b.id = m.brand_id
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Println("get brands error: ", err)
		return []Brand{}, err
	}

	defer rows.Close()

	brandMap := make(map[int64]*Brand)
	for rows.Next() {
		var brandID int64
		var brandName, fileName string
		var modelName sql.NullString
		var modelID, modelBrandID sql.NullInt64

		err := rows.Scan(&brandID, &brandName, &fileName, &modelID, &modelName, &modelBrandID)
		if err != nil {
			fmt.Println("getbrands error: ", err)
			return []Brand{}, err
		}

		brand, exists := brandMap[brandID]
		if !exists {
			brand = &Brand{
				ID:        brandID,
				BrandName: brandName,
				FileName:  fileName,
				Models:    []Model{},
			}
			brandMap[brandID] = brand
		}

		if modelID.Valid {
			brand.Models = append(brand.Models, Model{
				ID:        modelID.Int64,
				ModelName: modelName.String,
				BrandID:   modelBrandID.Int64,
			})
		}
	}

	var brands []Brand
	for _, brand := range brandMap {
		brands = append(brands, *brand)
	}

	sort.Slice(brands, func(i, j int) bool {
		return brands[i].ID < brands[j].ID
	})

	return brands, nil
}

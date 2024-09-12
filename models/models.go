package models

import (
	"example.com/db"
)

type Model struct {
	ID        int64
	ModelName string `binding:"required"`
	BrandID   int64  `binding:"required"`
}

func (m *Model) Save() (Model, error) {
	query := "INSERT INTO models(model_name, brand_id) VALUES ($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Model{}, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(m.ModelName, m.BrandID).Scan(&m.ID)
	if err != nil {
		return Model{}, err
	}

	return *m, nil
}

func GetModels() ([]Model, error) {
	query := "SELECT id, model_name, brand_id FROM models"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var models []Model
	for rows.Next() {
		var model Model
		err = rows.Scan(&model.ID, &model.ModelName, &model.BrandID)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

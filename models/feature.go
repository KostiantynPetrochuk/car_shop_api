package models

import (
	"example.com/db"
)

type Feature struct {
	ID          int64
	FeatureName string `binding:"required"`
	Category    string `binding:"required"`
}

func (f *Feature) Save() error {
	query := "INSERT INTO features(feature_name, category) VALUES ($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(f.FeatureName, f.Category).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetFeatures() ([]Feature, error) {
	query := "SELECT id, feature_name, category FROM features"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var features []Feature
	for rows.Next() {
		var feature Feature
		err := rows.Scan(&feature.ID, &feature.FeatureName, &feature.Category)
		if err != nil {
			return nil, err
		}

		features = append(features, feature)
	}

	return features, nil
}

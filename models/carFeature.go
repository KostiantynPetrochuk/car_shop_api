package models

import (
	"example.com/db"
)

type CarFeature struct {
	ID        int64
	CarId     int64
	FeatureId int64
}

func (cf *CarFeature) Save() error {
	query := "INSERT INTO car_features(car_id, feature_id) VALUES ($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(cf.CarId, cf.FeatureId).Scan(&cf.ID)
	if err != nil {
		return err
	}

	return nil
}

func SaveManyCarFeatures(carId int64, featureIds []int64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, featureId := range featureIds {
		carFeature := CarFeature{
			CarId:     carId,
			FeatureId: featureId,
		}

		if err := carFeature.Save(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

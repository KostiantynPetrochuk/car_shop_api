package routes

import (
	"net/http"

	"example.com/models"
	"github.com/gin-gonic/gin"
)

func addFeature(context *gin.Context) {
	var requestData struct {
		FeatureName string `json:"featureName"`
		Category    string `json:"category"`
	}

	if err := context.BindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		return
	}

	var feature = models.Feature{
		FeatureName: requestData.FeatureName,
		Category:    requestData.Category,
	}

	err := feature.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save feature."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"feature": feature})
}

func getFeatures(context *gin.Context) {
	features, err := models.GetFeatures()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch features."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"features": features})
}

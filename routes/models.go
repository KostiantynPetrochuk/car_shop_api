package routes

import (
	"fmt"
	"net/http"

	"example.com/models"
	"github.com/gin-gonic/gin"
)

func addModel(context *gin.Context) {
	var requestData struct {
		ModelName string `json:"modelName"`
		BrandID   int64  `json:"brandId"`
	}

	if err := context.BindJSON(&requestData); err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		return
	}

	var model = models.Model{
		ModelName: requestData.ModelName,
		BrandID:   requestData.BrandID,
	}

	createdModel, err := model.Save()
	fmt.Println("createdModel: ", createdModel)
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save model."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"model": createdModel})
}

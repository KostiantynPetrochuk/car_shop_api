package routes

import (
	"fmt"
	"log"
	"net/http"

	"example.com/models"
	"github.com/gin-gonic/gin"
)

func addBrand(context *gin.Context) {
	file, _ := context.FormFile("file")
	log.Println(file.Filename)
	brand_name, _ := context.GetPostForm("brand_name")
	fmt.Println(brand_name)

	dst := "uploads/brands/" + file.Filename
	context.SaveUploadedFile(file, dst)

	var brand = models.Brand{
		BrandName: brand_name,
		FileName:  file.Filename,
	}

	createdBrand, err := brand.Save()
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save brand."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"brand": createdBrand})
}

func getBrands(context *gin.Context) {
	brands, err := models.GetBrands()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get brands."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"brands": brands})
}

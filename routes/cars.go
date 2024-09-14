package routes

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/models"
	filesService "example.com/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func addCar(context *gin.Context) {
	var car models.Car

	car.VIN = context.PostForm("vin")
	brandId, err := strconv.Atoi(context.PostForm("brandId"))
	if err != nil {
		fmt.Println("Error converting brandId:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid brand ID."})
		return
	}
	car.BrandId = brandId
	modelId, err := strconv.Atoi(context.PostForm("modelId"))
	if err != nil {
		fmt.Println("Error converting modelId:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid model ID."})
		return
	}
	car.ModelId = modelId
	car.BodyType = context.PostForm("body")
	car.FuelType = context.PostForm("fuel_type")
	car.Year, err = strconv.Atoi(context.PostForm("year"))
	if err != nil {
		fmt.Println("Error converting year:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid year."})
		return
	}
	car.Transmission = context.PostForm("transmission")
	car.DriveType = context.PostForm("drive_type")
	car.Condition = context.PostForm("condition")
	car.EngineSize, err = strconv.ParseFloat(context.PostForm("engine_size"), 64)
	if err != nil {
		fmt.Println("Error converting engine_size:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid engine size."})
		return
	}
	car.DoorCount, err = strconv.Atoi(context.PostForm("door_count"))
	if err != nil {
		fmt.Println("Error converting door_count:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid door count."})
		return
	}
	car.CylinderCount, err = strconv.Atoi(context.PostForm("cylinder_count"))
	if err != nil {
		fmt.Println("Error converting cylinder_count:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid cylinder count."})
		return
	}
	car.Color = context.PostForm("color")
	car.Mileage, err = strconv.Atoi(context.PostForm("mileage"))
	if err != nil {
		fmt.Println("Error converting mileage:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid mileage."})
		return
	}

	fmt.Println("featureIds: ", context.PostForm("featureIds"))
	featureIds := strings.Trim(context.PostForm("featureIds"), "[]")
	featureIdsSlice := strings.Split(featureIds, ",")
	var featureIdsInt []int64
	for _, featureId := range featureIdsSlice {
		featureId = strings.TrimSpace(featureId)
		featureIdInt, err := strconv.Atoi(featureId)
		if err != nil {
			fmt.Println("Error converting featureId:", featureId, err)
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid feature ID: " + featureId})
			return
		}
		featureIdsInt = append(featureIdsInt, int64(featureIdInt))
	}

	form, err := context.MultipartForm()
	if err != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]
	var imagePaths []string

	for _, file := range files {
		firstFilename := filepath.Base(file.Filename)
		fileExt := filepath.Ext(firstFilename)
		randomFileName := uuid.New().String()
		filename := randomFileName + fileExt
		imagePath := fmt.Sprintf("uploads/cars/%s", filename)
		originalImagePath := imagePath
		if err := context.SaveUploadedFile(file, imagePath); err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		realFileName := strings.TrimSuffix(filename, filepath.Ext(filename))
		fmt.Println("realFileName: ", realFileName)

		if filepath.Ext(filename) == ".HEIC" {
			webpPath := fmt.Sprintf("uploads/cars/%s.webp", realFileName)
			if err := filesService.ConvertHeicToWebp(imagePath, webpPath); err != nil {
				context.String(http.StatusBadRequest, fmt.Sprintf("convert HEIC to JPEG err: %s", err.Error()))
				return
			}
			imagePath = randomFileName + ".webp"
			fmt.Println("originalImagePath: ", originalImagePath)
			if err := os.Remove(originalImagePath); err != nil {
				context.String(http.StatusBadRequest, fmt.Sprintf("remove original HEIC file err: %s", err.Error()))
				return
			}
		}
		if filepath.Ext(filename) != ".webp" && filepath.Ext(filename) != ".HEIC" {
			webpPath := fmt.Sprintf("uploads/cars/%s.webp", realFileName)
			if err := filesService.ConvertToWebp(imagePath, webpPath); err != nil {
				context.String(http.StatusBadRequest, fmt.Sprintf("convert to webp err: %s", err.Error()))
				return
			}
			imagePaths = append(imagePaths, randomFileName+".webp")
			if err := os.Remove(imagePath); err != nil {
				context.String(http.StatusBadRequest, fmt.Sprintf("remove original file err: %s", err.Error()))
				return
			}
		} else {
			imagePaths = append(imagePaths, randomFileName+".webp")
		}
	}

	imagePathsJSON, err := json.Marshal(imagePaths)
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car."})
		return
	}
	car.ImageNames = json.RawMessage(imagePathsJSON)

	if err := car.Save(); err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car."})
		return
	}

	fmt.Println("featureIdsInt: ", featureIdsInt)

	if err := models.SaveManyCarFeatures(car.ID, featureIdsInt); err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car features."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getCars(context *gin.Context) {
	cars, err := models.GetCars()
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get cars."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"cars": cars})
}

func getCar(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID."})
		return
	}
	car, err := models.GetCarById(int64(id))
	if err != nil {
		fmt.Println("error: ", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"car": car})
}

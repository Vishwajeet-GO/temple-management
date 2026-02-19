package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

func CreateFestival(festival *models.Festival) error {
	return database.DB.Create(festival).Error
}

func GetAllFestivals() ([]models.Festival, error) {
	var festivals []models.Festival
	err := database.DB.Order("created_at DESC").Find(&festivals).Error
	return festivals, err
}

func GetFestivalByID(id uint) (*models.Festival, error) {
	var festival models.Festival
	err := database.DB.First(&festival, id).Error
	return &festival, err
}

func UpdateFestival(festival *models.Festival) error {
	return database.DB.Save(festival).Error
}

func DeleteFestival(id uint) error {
	return database.DB.Delete(&models.Festival{}, id).Error
}

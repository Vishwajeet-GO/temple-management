package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

func GetTempleInfo() (*models.TempleInfo, error) {
	var info models.TempleInfo
	err := database.DB.First(&info).Error
	return &info, err
}

func UpdateTempleInfo(info *models.TempleInfo) error {
	return database.DB.Save(info).Error
}

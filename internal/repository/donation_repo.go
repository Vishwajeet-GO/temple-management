package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

func CreateDonation(donation *models.Donation) error {
	return database.DB.Create(donation).Error
}

func GetAllDonations() ([]models.Donation, error) {
	var donations []models.Donation
	err := database.DB.Preload("Festival").Order("created_at DESC").Find(&donations).Error
	return donations, err
}

func GetDonationByID(id uint) (*models.Donation, error) {
	var donation models.Donation
	err := database.DB.Preload("Festival").First(&donation, id).Error
	return &donation, err
}

func UpdateDonation(donation *models.Donation) error {
	return database.DB.Save(donation).Error
}

func DeleteDonation(id uint) error {
	return database.DB.Delete(&models.Donation{}, id).Error
}

func GetDonationsByFestival(festivalID uint) ([]models.Donation, error) {
	var donations []models.Donation
	err := database.DB.Where("festival_id = ?", festivalID).Order("created_at DESC").Find(&donations).Error
	return donations, err
}

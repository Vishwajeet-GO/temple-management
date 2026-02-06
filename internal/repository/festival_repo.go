package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

type FestivalRepository struct{}

func NewFestivalRepository() *FestivalRepository {
	return &FestivalRepository{}
}

func (r *FestivalRepository) Create(festival *models.Festival) error {
	return database.GetDB().Create(festival).Error
}

func (r *FestivalRepository) GetAll() ([]models.Festival, error) {
	var festivals []models.Festival
	result := database.GetDB().Order("start_date desc").Find(&festivals)
	return festivals, result.Error
}

func (r *FestivalRepository) GetByID(id uint) (*models.Festival, error) {
	var festival models.Festival
	result := database.GetDB().Preload("Donations").Preload("Expenses").First(&festival, id)
	return &festival, result.Error
}

func (r *FestivalRepository) Update(festival *models.Festival) error {
	return database.GetDB().Save(festival).Error
}

func (r *FestivalRepository) Delete(id uint) error {
	return database.GetDB().Delete(&models.Festival{}, id).Error
}

func (r *FestivalRepository) GetUpcoming() ([]models.Festival, error) {
	var festivals []models.Festival
	result := database.GetDB().Where("status = ?", "upcoming").Order("start_date asc").Find(&festivals)
	return festivals, result.Error
}

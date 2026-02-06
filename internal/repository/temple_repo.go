package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

type TempleRepository struct{}

func NewTempleRepository() *TempleRepository {
	return &TempleRepository{}
}

func (r *TempleRepository) GetTemple() (*models.Temple, error) {
	var temple models.Temple
	result := database.GetDB().First(&temple)
	return &temple, result.Error
}

func (r *TempleRepository) UpdateTemple(temple *models.Temple) error {
	return database.GetDB().Save(temple).Error
}

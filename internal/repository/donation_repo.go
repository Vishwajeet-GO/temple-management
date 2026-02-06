package repository

import (
	"fmt"
	"newapp/internal/database"
	"newapp/internal/models"
	"time"
)

type DonationRepository struct{}

func NewDonationRepository() *DonationRepository {
	return &DonationRepository{}
}

func (r *DonationRepository) Create(donation *models.Donation) error {
	// Generate receipt number
	donation.ReceiptNo = fmt.Sprintf("RCP-%d-%d", time.Now().Unix(), donation.TempleID)
	if donation.Date.IsZero() {
		donation.Date = time.Now()
	}
	return database.GetDB().Create(donation).Error
}

func (r *DonationRepository) GetAll() ([]models.Donation, error) {
	var donations []models.Donation
	result := database.GetDB().Preload("Festival").Order("date desc").Find(&donations)
	return donations, result.Error
}

func (r *DonationRepository) GetByID(id uint) (*models.Donation, error) {
	var donation models.Donation
	result := database.GetDB().Preload("Festival").First(&donation, id)
	return &donation, result.Error
}

func (r *DonationRepository) GetByFestival(festivalID uint) ([]models.Donation, error) {
	var donations []models.Donation
	result := database.GetDB().Where("festival_id = ?", festivalID).Order("date desc").Find(&donations)
	return donations, result.Error
}

func (r *DonationRepository) Update(donation *models.Donation) error {
	return database.GetDB().Save(donation).Error
}

func (r *DonationRepository) Delete(id uint) error {
	return database.GetDB().Delete(&models.Donation{}, id).Error
}

func (r *DonationRepository) GetTotal() (float64, error) {
	var total float64
	result := database.GetDB().Model(&models.Donation{}).Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total, result.Error
}

func (r *DonationRepository) GetTotalByFestival(festivalID uint) (float64, error) {
	var total float64
	result := database.GetDB().Model(&models.Donation{}).Where("festival_id = ?", festivalID).Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total, result.Error
}

func (r *DonationRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Donation, error) {
	var donations []models.Donation
	result := database.GetDB().Where("date BETWEEN ? AND ?", startDate, endDate).Order("date desc").Find(&donations)
	return donations, result.Error
}

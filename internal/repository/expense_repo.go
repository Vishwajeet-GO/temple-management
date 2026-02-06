package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
	"time"
)

type ExpenseRepository struct{}

func NewExpenseRepository() *ExpenseRepository {
	return &ExpenseRepository{}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}
	return database.GetDB().Create(expense).Error
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	var expenses []models.Expense
	result := database.GetDB().Preload("Festival").Order("date desc").Find(&expenses)
	return expenses, result.Error
}

func (r *ExpenseRepository) GetByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	result := database.GetDB().Preload("Festival").First(&expense, id)
	return &expense, result.Error
}

func (r *ExpenseRepository) GetByFestival(festivalID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := database.GetDB().Where("festival_id = ?", festivalID).Order("date desc").Find(&expenses)
	return expenses, result.Error
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return database.GetDB().Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uint) error {
	return database.GetDB().Delete(&models.Expense{}, id).Error
}

func (r *ExpenseRepository) GetTotal() (float64, error) {
	var total float64
	result := database.GetDB().Model(&models.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total, result.Error
}

func (r *ExpenseRepository) GetTotalByFestival(festivalID uint) (float64, error) {
	var total float64
	result := database.GetDB().Model(&models.Expense{}).Where("festival_id = ?", festivalID).Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total, result.Error
}

func (r *ExpenseRepository) GetByCategory(category string) ([]models.Expense, error) {
	var expenses []models.Expense
	result := database.GetDB().Where("category = ?", category).Order("date desc").Find(&expenses)
	return expenses, result.Error
}

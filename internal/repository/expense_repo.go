package repository

import (
	"newapp/internal/database"
	"newapp/internal/models"
)

func CreateExpense(expense *models.Expense) error {
	return database.DB.Create(expense).Error
}

func GetAllExpenses() ([]models.Expense, error) {
	var expenses []models.Expense
	err := database.DB.Preload("Festival").Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

func GetExpenseByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	err := database.DB.Preload("Festival").First(&expense, id).Error
	return &expense, err
}

func UpdateExpense(expense *models.Expense) error {
	return database.DB.Save(expense).Error
}

func DeleteExpense(id uint) error {
	return database.DB.Delete(&models.Expense{}, id).Error
}

func GetExpensesByFestival(festivalID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	err := database.DB.Where("festival_id = ?", festivalID).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

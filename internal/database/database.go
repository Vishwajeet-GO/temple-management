package database

import (
	"fmt"
	"log"
	"newapp/internal/config"
	"newapp/internal/models"
	"os"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) {
	var err error
	var dialector gorm.Dialector

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		log.Println("Using PostgreSQL from DATABASE_URL")
		dialector = postgres.Open(databaseURL)
	} else if cfg.DBType == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
		dialector = postgres.Open(dsn)
	} else {
		log.Println("Using SQLite for local development")
		dialector = sqlite.Open(cfg.DBName)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = DB.AutoMigrate(
		&models.User{},
		&models.Temple{},
		&models.Festival{},
		&models.Donation{},
		&models.Expense{},
		&models.Event{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedInitialData()
}

func seedInitialData() {
	// 1. Admin User Setup
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Name:     "Temple Admin",
			Role:     "admin",
			IsActive: true,
		}
		DB.Create(&adminUser)
		log.Println("Default admin user created")
	}

	// 2. TEMPLE DETAILS UPDATE (English -> Hindi)
	var temple models.Temple
	// Check if temple exists
	result := DB.First(&temple)

	// 👇 CHANGE YOUR DETAILS HERE 👇
	newTempleName := "श्री गौरी शंकर मंदिर "                             // Change to your Hindi Name
	newAddress := "Poisar, Kandivali East, Mumbai, Maharashtra - 400101" // Change to Full Address
	newCity := "Mumbai"
	newState := "Maharashtra"
	// 👆 CHANGE END 👆

	if result.Error != nil {
		// If no temple exists, create it
		temple = models.Temple{
			Name:            newTempleName,
			Address:         newAddress,
			City:            newCity,
			State:           newState,
			PinCode:         "400001",
			Phone:           "9876543210",
			Email:           "temple@example.com",
			Description:     "Jay Shree Ram",
			MainDeity:       "Hanuman Ji",
			EstablishedYear: 1950,
		}
		DB.Create(&temple)
		log.Println("Initial temple data created!")
	} else {
		// If temple exists, UPDATE it
		temple.Name = newTempleName
		temple.Address = newAddress
		temple.City = newCity
		temple.State = newState
		DB.Save(&temple)
		log.Println("Temple details updated to Hindi!")
	}
}

func GetDB() *gorm.DB {
	return DB
}

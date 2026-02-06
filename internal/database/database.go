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

	// Check for DATABASE_URL first (Render provides this)
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		// Use PostgreSQL from DATABASE_URL
		log.Println("Using PostgreSQL from DATABASE_URL")
		dialector = postgres.Open(databaseURL)
	} else if cfg.DBType == "postgres" {
		// Use PostgreSQL from individual config
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
		log.Println("Using PostgreSQL from config")
		dialector = postgres.Open(dsn)
	} else {
		// Use SQLite for local development
		log.Println("Using SQLite for local development")
		dialector = sqlite.Open(cfg.DBName)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto migrate the schemas
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

	log.Println("Database migrated successfully!")

	// Seed initial data
	seedInitialData()
}

func seedInitialData() {
	// Create default admin user if not exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1122"), bcrypt.DefaultCost)

		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Name:     "Temple Admin",
			Role:     "admin",
			IsActive: true,
		}
		DB.Create(&adminUser)
		log.Println("========================================")
		log.Println("Default admin user created!")
		log.Println("Username: admin")
		log.Println("Password: admin123")
		log.Println("WARNING: Change this password after first login!")
		log.Println("========================================")
	}

	// Create default temple if not exists
	var templeCount int64
	DB.Model(&models.Temple{}).Count(&templeCount)

	if templeCount == 0 {
		temple := models.Temple{
			Name:            "Shri Ram Mandir",
			Address:         "Main Temple Road",
			City:            "Your City",
			State:           "Your State",
			PinCode:         "123456",
			Phone:           "9876543210",
			Email:           "temple@example.com",
			Description:     "A beautiful temple dedicated to Lord Ram",
			MainDeity:       "Lord Ram",
			EstablishedYear: 1950,
		}
		DB.Create(&temple)
		log.Println("Initial temple data seeded!")
	}
}

func GetDB() *gorm.DB {
	return DB
}

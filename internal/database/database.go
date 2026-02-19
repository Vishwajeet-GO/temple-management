package database

import (
	"fmt"
	"log"
	"time"

	"newapp/internal/config"
	"newapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Kolkata",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	log.Printf("🔌 Connecting to PostgreSQL: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Database connection failed: ", err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ Database ping failed: ", err)
	}

	log.Println("✅ PostgreSQL connected")

	// Migrate
	DB.AutoMigrate(
		&models.User{},
		&models.Festival{},
		&models.Donation{},
		&models.Expense{},
		&models.TempleInfo{},
	)
	log.Println("✅ Tables migrated")

	// Default admin
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		DB.Create(&models.User{Username: "admin", Password: "admin123", Role: "admin"})
		log.Println("✅ Default admin: admin / admin123")
	}

	// Default temple
	DB.Model(&models.TempleInfo{}).Count(&count)
	if count == 0 {
		DB.Create(&models.TempleInfo{Name: "Temple Management", UPI: "8097890684@mbk"})
	}
}

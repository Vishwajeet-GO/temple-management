package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"newapp/internal/config"
	"newapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(cfg *config.Config) {
	var dsn string

	// Priority 1: DATABASE_URL (Render provides this)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
		log.Println("🔌 Connecting using DATABASE_URL")
	} else if cfg.DBHost != "localhost" {
		// Priority 2: Individual env vars (non-localhost = production)
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require TimeZone=Asia/Kolkata",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
		log.Printf("🔌 Connecting to PostgreSQL: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	} else {
		// Priority 3: Localhost (development)
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Kolkata",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
		log.Printf("🔌 Connecting to local PostgreSQL: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

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

	// Force reset admin password (remove this after first deploy)
	// var adminUser models.User
	// if DB.Where("username = ?", "admin").First(&adminUser).Error == nil {
	// 	DB.Model(&adminUser).Update("password", "admin123")
	// 	log.Println("✅ Admin password reset to: admin123")
	// } else {
	// 	DB.Create(&models.User{Username: "admin", Password: "admin123", Role: "admin"})
	// 	log.Println("✅ Default admin created: admin / admin123")
	// }

	// Default temple
	var templeCount int64
	DB.Model(&models.TempleInfo{}).Count(&templeCount)
	if templeCount == 0 {
		DB.Create(&models.TempleInfo{
			Name:  "श्री गौरी शंकर मंदिर",
			UPI:   "8097890684@mbk",
			Phone: "8097890684",
		})
	}
}

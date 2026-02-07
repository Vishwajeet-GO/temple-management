package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Name     string `json:"name"`
	Role     string `json:"role" gorm:"default:'admin'"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

type Temple struct {
	gorm.Model
	Name            string `json:"name"`
	Address         string `json:"address"`
	City            string `json:"city"`
	State           string `json:"state"`
	PinCode         string `json:"pin_code"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Description     string `json:"description"`
	MainDeity       string `json:"main_deity"`
	EstablishedYear int    `json:"established_year"`
}

type Festival struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	TempleID    uint      `json:"temple_id"`
}

type Donation struct {
	gorm.Model
	DonorName      string    `json:"donor_name"`
	DonorPhone     string    `json:"donor_phone"`
	DonorAddress   string    `json:"donor_address"`
	Amount         float64   `json:"amount"`
	PaymentMode    string    `json:"payment_mode"`
	Purpose        string    `json:"purpose"`
	ReceiptNo      string    `json:"receipt_no"`
	Date           time.Time `json:"date"`
	Notes          string    `json:"notes"`
	ScreenshotPath string    `json:"screenshot_path"`
	TempleID       uint      `json:"temple_id"`
	FestivalID     uint      `json:"festival_id"`
}

type Expense struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	PaymentMode string    `json:"payment_mode"`
	VendorName  string    `json:"vendor_name"`
	BillNo      string    `json:"bill_no"`
	Date        time.Time `json:"date"`
	TempleID    uint      `json:"temple_id"`
	FestivalID  uint      `json:"festival_id"`
}

type Event struct {
	gorm.Model
	// Add fields if needed later
}

type DashboardSummary struct {
	TotalDonations    float64 `json:"total_donations"`
	TotalExpenses     float64 `json:"total_expenses"`
	Balance           float64 `json:"balance"`
	TotalFestivals    int64   `json:"total_festivals"`
	UpcomingFestivals int64   `json:"upcoming_festivals"`
	TotalEvents       int64   `json:"total_events"` // Added
	DonorsCount       int64   `json:"donors_count"` // Added
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"` // Added
}

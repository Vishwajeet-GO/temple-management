package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents admin users who can login
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Name     string `json:"name"`
	Role     string `json:"role" gorm:"default:'admin'"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

// Temple represents the main temple information
type Temple struct {
	gorm.Model
	Name            string `json:"name" gorm:"not null"`
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

// Festival represents temple festivals/events
type Festival struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Status      string     `json:"status" gorm:"default:'upcoming'"`
	TempleID    uint       `json:"temple_id"`
	Temple      Temple     `json:"temple" gorm:"foreignKey:TempleID"`
	Donations   []Donation `json:"donations" gorm:"foreignKey:FestivalID"`
	Expenses    []Expense  `json:"expenses" gorm:"foreignKey:FestivalID"`
}

// Donation represents money collected
type Donation struct {
	gorm.Model
	DonorName    string    `json:"donor_name"`
	DonorPhone   string    `json:"donor_phone"`
	DonorAddress string    `json:"donor_address"`
	Amount       float64   `json:"amount" gorm:"not null"`
	PaymentMode  string    `json:"payment_mode"`
	Purpose      string    `json:"purpose"`
	ReceiptNo    string    `json:"receipt_no" gorm:"unique"`
	Date         time.Time `json:"date"`
	FestivalID   *uint     `json:"festival_id"`
	Festival     *Festival `json:"festival" gorm:"foreignKey:FestivalID"`
	TempleID     uint      `json:"temple_id"`
	Temple       Temple    `json:"temple" gorm:"foreignKey:TempleID"`
	Notes        string    `json:"notes"`
}

// Expense represents money spent
type Expense struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Category    string    `json:"category"`
	PaymentMode string    `json:"payment_mode"`
	VendorName  string    `json:"vendor_name"`
	BillNo      string    `json:"bill_no"`
	Date        time.Time `json:"date"`
	FestivalID  *uint     `json:"festival_id"`
	Festival    *Festival `json:"festival" gorm:"foreignKey:FestivalID"`
	TempleID    uint      `json:"temple_id"`
	Temple      Temple    `json:"temple" gorm:"foreignKey:TempleID"`
	ApprovedBy  string    `json:"approved_by"`
}

// Event represents general temple events
type Event struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Venue       string    `json:"venue"`
	Organizer   string    `json:"organizer"`
	TempleID    uint      `json:"temple_id"`
	Temple      Temple    `json:"temple" gorm:"foreignKey:TempleID"`
}

// DashboardSummary for dashboard
type DashboardSummary struct {
	TotalDonations    float64 `json:"total_donations"`
	TotalExpenses     float64 `json:"total_expenses"`
	Balance           float64 `json:"balance"`
	TotalFestivals    int64   `json:"total_festivals"`
	UpcomingFestivals int64   `json:"upcoming_festivals"`
	TotalEvents       int64   `json:"total_events"`
	DonorsCount       int64   `json:"donors_count"`
}

// LoginRequest for login API
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse for login API response
type LoginResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

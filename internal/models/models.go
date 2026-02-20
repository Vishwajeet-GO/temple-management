package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:viewer"`
	CreatedAt time.Time `json:"created_at"`
}

type Festival struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:upcoming"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Donation struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FestivalID  *uint     `json:"festival_id" gorm:"index"`
	Festival    *Festival `json:"festival,omitempty" gorm:"foreignKey:FestivalID"`
	Date        string    `json:"date"`
	Donor       string    `json:"donor"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status" gorm:"default:pending"`
	PaymentMode string    `json:"payment_mode" gorm:"default:cash"`
	ImageURL    string    `json:"image_url"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Expense struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FestivalID *uint     `json:"festival_id" gorm:"index"`
	Festival   *Festival `json:"festival,omitempty" gorm:"foreignKey:FestivalID"`
	Title      string    `json:"title"`
	Amount     float64   `json:"amount"`
	Date       string    `json:"date"`
	Status     string    `json:"status" gorm:"default:pending"`
	Category   string    `json:"category"`
	Note       string    `json:"note"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TempleInfo struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	UPI     string `json:"upi"`
	QRCode  string `json:"qr_code"`
}

type GalleryItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Category  string    `json:"category" gorm:"default:temple"`
	Type      string    `json:"type" gorm:"default:photo"` // photo or video
	URL       string    `json:"url"`
	Thumbnail string    `json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`
}

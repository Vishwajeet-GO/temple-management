package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubmitDonation handles public donation from the donate page
func SubmitDonation(c *gin.Context) {
	var donation models.Donation
	ct := c.ContentType()

	if strings.Contains(ct, "multipart") {
		donation.Donor = c.PostForm("donor")
		donation.PaymentMode = c.PostForm("payment_mode")
		donation.Link = c.PostForm("link")
		donation.Date = time.Now().Format("2006-01-02")
		donation.Status = "pending"

		amt, err := strconv.ParseFloat(c.PostForm("amount"), 64)
		if err == nil {
			donation.Amount = amt
		}

		// Handle festival_id
		if fid := c.PostForm("festival_id"); fid != "" && fid != "0" {
			id, _ := strconv.ParseUint(fid, 10, 32)
			uid := uint(id)
			donation.FestivalID = &uid
		}

		// Handle screenshot upload
		file, err := c.FormFile("screenshot")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
			if allowedExts[strings.ToLower(ext)] {
				fname := fmt.Sprintf("don_%s%s", uuid.New().String()[:8], ext)
				savePath := filepath.Join("uploads", "donations", fname)
				if c.SaveUploadedFile(file, savePath) == nil {
					donation.ImageURL = "/" + filepath.ToSlash(savePath)
				}
			}
		}
	} else {
		if err := c.ShouldBindJSON(&donation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid data: " + err.Error()})
			return
		}
		donation.Date = time.Now().Format("2006-01-02")
		donation.Status = "pending"
	}

	// Validate
	if donation.Donor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Donor name is required"})
		return
	}
	if donation.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Amount must be greater than 0"})
		return
	}
	if donation.PaymentMode == "" {
		donation.PaymentMode = "cash"
	}

	// Save
	if err := database.DB.Create(&donation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to save donation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Thank you for your donation! 🙏",
		"data":    donation,
	})
}

// GetPaymentInfo returns temple payment details (UPI, QR etc)
func GetPaymentInfo(c *gin.Context) {
	var info models.TempleInfo
	database.DB.First(&info)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":    info.Name,
			"upi":     info.UPI,
			"phone":   info.Phone,
			"qr_code": info.QRCode,
		},
	})
}

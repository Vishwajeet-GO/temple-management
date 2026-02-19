package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetDonations(c *gin.Context) {
	var donations []models.Donation
	search := c.Query("search")
	festivalID := c.Query("festival_id")

	query := database.DB.Preload("Festival").Order("created_at DESC")

	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		query = query.Where("festival_id = ?", festivalID)
	}
	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(donor) LIKE ? OR LOWER(payment_mode) LIKE ? OR CAST(amount AS TEXT) LIKE ?", s, s, s)
	}
	query.Find(&donations)

	// Stats - respect festival filter
	statQuery := database.DB.Model(&models.Donation{})
	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		statQuery = statQuery.Where("festival_id = ?", festivalID)
	}

	var total float64
	var paid, pending int64
	statQuery.Select("COALESCE(SUM(amount),0)").Scan(&total)
	sq2 := database.DB.Model(&models.Donation{})
	sq3 := database.DB.Model(&models.Donation{})
	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		sq2 = sq2.Where("festival_id = ?", festivalID)
		sq3 = sq3.Where("festival_id = ?", festivalID)
	}
	sq2.Where("status=?", "paid").Count(&paid)
	sq3.Where("status=?", "pending").Count(&pending)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    donations,
		"stats": gin.H{
			"total":   total,
			"entries": len(donations),
			"paid":    paid,
			"pending": pending,
		},
	})
}

func CreateDonation(c *gin.Context) {
	var donation models.Donation
	ct := c.ContentType()

	if strings.Contains(ct, "multipart") {
		donation.Date = c.PostForm("date")
		donation.Donor = c.PostForm("donor")
		amt, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
		donation.Amount = amt
		donation.Status = c.PostForm("status")
		donation.PaymentMode = c.PostForm("payment_mode")
		donation.Link = c.PostForm("link")

		if fid := c.PostForm("festival_id"); fid != "" && fid != "0" {
			id, _ := strconv.ParseUint(fid, 10, 32)
			uid := uint(id)
			donation.FestivalID = &uid
		}

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			fname := fmt.Sprintf("don_%s%s", uuid.New().String()[:8], ext)
			savePath := filepath.Join("uploads", "donations", fname)
			os.MkdirAll(filepath.Dir(savePath), 0755)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				donation.ImageURL = "/" + filepath.ToSlash(savePath)
			}
		}
	} else {
		if err := c.ShouldBindJSON(&donation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid data: " + err.Error()})
			return
		}
	}

	if donation.Date == "" {
		donation.Date = time.Now().Format("2006-01-02")
	}
	if donation.Status == "" {
		donation.Status = "pending"
	}
	if donation.PaymentMode == "" {
		donation.PaymentMode = "cash"
	}
	if donation.Donor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Donor name required"})
		return
	}
	if donation.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Amount must be > 0"})
		return
	}

	database.DB.Create(&donation)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": donation, "message": "Donation added"})
}

func UpdateDonation(c *gin.Context) {
	id := c.Param("id")
	var donation models.Donation
	if database.DB.First(&donation, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}

	ct := c.ContentType()
	if strings.Contains(ct, "multipart") {
		if v := c.PostForm("donor"); v != "" {
			donation.Donor = v
		}
		if v := c.PostForm("amount"); v != "" {
			a, _ := strconv.ParseFloat(v, 64)
			donation.Amount = a
		}
		if v := c.PostForm("date"); v != "" {
			donation.Date = v
		}
		if v := c.PostForm("status"); v != "" {
			donation.Status = v
		}
		if v := c.PostForm("payment_mode"); v != "" {
			donation.PaymentMode = v
		}
		if fid := c.PostForm("festival_id"); fid != "" && fid != "0" {
			id, _ := strconv.ParseUint(fid, 10, 32)
			uid := uint(id)
			donation.FestivalID = &uid
		}
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			fname := fmt.Sprintf("don_%s%s", uuid.New().String()[:8], ext)
			savePath := filepath.Join("uploads", "donations", fname)
			os.MkdirAll(filepath.Dir(savePath), 0755)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				if donation.ImageURL != "" {
					os.Remove("." + donation.ImageURL)
				}
				donation.ImageURL = "/" + filepath.ToSlash(savePath)
			}
		}
	} else {
		var u models.Donation
		c.ShouldBindJSON(&u)
		database.DB.Model(&donation).Updates(u)
	}

	database.DB.Save(&donation)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": donation, "message": "Updated"})
}

func DeleteDonation(c *gin.Context) {
	id := c.Param("id")
	var d models.Donation
	if database.DB.First(&d, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}
	if d.ImageURL != "" {
		os.Remove("." + d.ImageURL)
	}
	database.DB.Delete(&d)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted"})
}

func ToggleDonation(c *gin.Context) {
	id := c.Param("id")
	var d models.Donation
	if database.DB.First(&d, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}
	ns := "paid"
	if d.Status == "paid" {
		ns = "pending"
	}
	database.DB.Model(&d).Update("status", ns)
	d.Status = ns
	c.JSON(http.StatusOK, gin.H{"success": true, "data": d, "message": "Status: " + ns})
}

func SubmitDonations(c *gin.Context) {
	var d models.Donation
	ct := c.ContentType()
	if strings.Contains(ct, "multipart") {
		d.Donor = c.PostForm("donor")
		amt, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
		d.Amount = amt
		d.PaymentMode = c.PostForm("payment_mode")
		d.Date = time.Now().Format("2006-01-02")
		d.Status = "pending"
		file, err := c.FormFile("screenshot")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			fname := fmt.Sprintf("don_%s%s", uuid.New().String()[:8], ext)
			savePath := filepath.Join("uploads", "donations", fname)
			os.MkdirAll(filepath.Dir(savePath), 0755)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				d.ImageURL = "/" + filepath.ToSlash(savePath)
			}
		}
	} else {
		c.ShouldBindJSON(&d)
		d.Date = time.Now().Format("2006-01-02")
		d.Status = "pending"
	}
	if d.Donor == "" || d.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Name and amount required"})
		return
	}
	database.DB.Create(&d)
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Thank you!"})
}

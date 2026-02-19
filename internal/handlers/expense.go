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

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	search := c.Query("search")
	festivalID := c.Query("festival_id")

	query := database.DB.Preload("Festival").Order("created_at DESC")

	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		query = query.Where("festival_id = ?", festivalID)
	}
	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(category) LIKE ? OR CAST(amount AS TEXT) LIKE ?", s, s, s)
	}
	query.Find(&expenses)

	statQuery := database.DB.Model(&models.Expense{})
	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		statQuery = statQuery.Where("festival_id = ?", festivalID)
	}
	var total float64
	var paid, pending int64
	statQuery.Select("COALESCE(SUM(amount),0)").Scan(&total)
	sq2 := database.DB.Model(&models.Expense{})
	sq3 := database.DB.Model(&models.Expense{})
	if festivalID != "" && festivalID != "0" && festivalID != "all" {
		sq2 = sq2.Where("festival_id = ?", festivalID)
		sq3 = sq3.Where("festival_id = ?", festivalID)
	}
	sq2.Where("status=?", "paid").Count(&paid)
	sq3.Where("status=?", "pending").Count(&pending)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    expenses,
		"stats":   gin.H{"total": total, "entries": len(expenses), "paid": paid, "pending": pending},
	})
}

func CreateExpense(c *gin.Context) {
	var expense models.Expense
	ct := c.ContentType()

	if strings.Contains(ct, "multipart") {
		expense.Title = c.PostForm("title")
		amt, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
		expense.Amount = amt
		expense.Date = c.PostForm("date")
		expense.Status = c.PostForm("status")
		expense.Category = c.PostForm("category")
		expense.Note = c.PostForm("note")

		if fid := c.PostForm("festival_id"); fid != "" && fid != "0" {
			id, _ := strconv.ParseUint(fid, 10, 32)
			uid := uint(id)
			expense.FestivalID = &uid
		}

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			fname := fmt.Sprintf("exp_%s%s", uuid.New().String()[:8], ext)
			savePath := filepath.Join("uploads", "expenses", fname)
			os.MkdirAll(filepath.Dir(savePath), 0755)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				expense.ImageURL = "/" + filepath.ToSlash(savePath)
			}
		}
	} else {
		if err := c.ShouldBindJSON(&expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid data: " + err.Error()})
			return
		}
	}

	if expense.Date == "" {
		expense.Date = time.Now().Format("2006-01-02")
	}
	if expense.Status == "" {
		expense.Status = "pending"
	}
	if expense.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Title required"})
		return
	}

	database.DB.Create(&expense)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": expense, "message": "Expense added"})
}

func UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	var expense models.Expense
	if database.DB.First(&expense, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}

	ct := c.ContentType()
	if strings.Contains(ct, "multipart") {
		if v := c.PostForm("title"); v != "" {
			expense.Title = v
		}
		if v := c.PostForm("amount"); v != "" {
			a, _ := strconv.ParseFloat(v, 64)
			expense.Amount = a
		}
		if v := c.PostForm("date"); v != "" {
			expense.Date = v
		}
		if v := c.PostForm("status"); v != "" {
			expense.Status = v
		}
		if v := c.PostForm("category"); v != "" {
			expense.Category = v
		}
		if v := c.PostForm("note"); v != "" {
			expense.Note = v
		}
		if fid := c.PostForm("festival_id"); fid != "" && fid != "0" {
			fuid, _ := strconv.ParseUint(fid, 10, 32)
			uid := uint(fuid)
			expense.FestivalID = &uid
		}
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			fname := fmt.Sprintf("exp_%s%s", uuid.New().String()[:8], ext)
			savePath := filepath.Join("uploads", "expenses", fname)
			os.MkdirAll(filepath.Dir(savePath), 0755)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				if expense.ImageURL != "" {
					os.Remove("." + expense.ImageURL)
				}
				expense.ImageURL = "/" + filepath.ToSlash(savePath)
			}
		}
	} else {
		var u models.Expense
		c.ShouldBindJSON(&u)
		database.DB.Model(&expense).Updates(u)
	}

	database.DB.Save(&expense)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": expense, "message": "Updated"})
}

func DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	var e models.Expense
	if database.DB.First(&e, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}
	if e.ImageURL != "" {
		os.Remove("." + e.ImageURL)
	}
	database.DB.Delete(&e)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted"})
}

func ToggleExpense(c *gin.Context) {
	id := c.Param("id")
	var e models.Expense
	if database.DB.First(&e, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
		return
	}
	ns := "paid"
	if e.Status == "paid" {
		ns = "pending"
	}
	database.DB.Model(&e).Update("status", ns)
	e.Status = ns
	c.JSON(http.StatusOK, gin.H{"success": true, "data": e, "message": "Status: " + ns})
}

package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGallery(c *gin.Context) {
	var items []models.GalleryItem
	category := c.Query("category")
	itemType := c.Query("type")

	query := database.DB.Order("created_at DESC")

	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}
	if itemType != "" && itemType != "all" {
		query = query.Where("type = ?", itemType)
	}

	query.Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
		"total":   len(items),
	})
}

func UploadGallery(c *gin.Context) {
	title := c.PostForm("title")
	category := c.PostForm("category")
	itemType := c.PostForm("type")

	if category == "" {
		category = "temple"
	}
	if itemType == "" {
		itemType = "photo"
	}
	if title == "" {
		title = "Untitled"
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "File is required"})
		return
	}

	// Validate file size (2048MB max for videos, 50MB for photos)
	maxSize := int64(50 * 1024 * 1024)
	if itemType == "video" {
		maxSize = 2048 * 1024 * 1024
	}
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "File too large"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedPhoto := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	allowedVideo := map[string]bool{".mp4": true, ".webm": true, ".mov": true}

	if itemType == "photo" && !allowedPhoto[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid image format. Use JPG, PNG, GIF or WebP"})
		return
	}
	if itemType == "video" && !allowedVideo[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid video format. Use MP4, WebM or MOV"})
		return
	}

	// Save file
	os.MkdirAll("uploads/gallery", 0755)
	fname := fmt.Sprintf("gal_%s%s", uuid.New().String()[:8], ext)
	savePath := filepath.Join("uploads", "gallery", fname)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to save file"})
		return
	}

	fileURL := "/" + filepath.ToSlash(savePath)

	item := models.GalleryItem{
		Title:     title,
		Category:  category,
		Type:      itemType,
		URL:       fileURL,
		Thumbnail: fileURL,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    item,
		"message": "Uploaded successfully!",
	})
}

func DeleteGalleryItem(c *gin.Context) {
	id := c.Param("id")
	var item models.GalleryItem

	if database.DB.First(&item, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Item not found"})
		return
	}

	// Delete file
	if item.URL != "" {
		os.Remove("." + item.URL)
	}

	database.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}

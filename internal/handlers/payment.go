package handlers

import (
    "fmt"
    "net/http"
    "path/filepath"
    "strconv"
    "time"

    "newapp/internal/database"
    "newapp/internal/models"

    "github.com/gin-gonic/gin"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
    return &PaymentHandler{}
}

// Process Donation with File Upload
func (h *PaymentHandler) ProcessDonation(c *gin.Context) {
    // 1. Get Form Fields
    name := c.PostForm("donor_name")
    phone := c.PostForm("donor_phone")
    amountStr := c.PostForm("amount")
    purpose := c.PostForm("purpose")
    paymentApp := c.PostForm("payment_app") // gpay, paytm, etc.

    if name == "" || amountStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Amount are required"})
        return
    }

    amount, _ := strconv.ParseFloat(amountStr, 64)
    receiptNo := fmt.Sprintf("RCP-%d", time.Now().Unix())

    // 2. Handle File Upload (Screenshot)
    var screenshotPath string
    file, err := c.FormFile("screenshot")
    if err == nil {
        // Save file
        filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
        uploadPath := "web/static/uploads/" + filename
        if err := c.SaveUploadedFile(file, uploadPath); err == nil {
            screenshotPath = "/static/uploads/" + filename
        }
    }

    // 3. Save to Database
    donation := models.Donation{
        DonorName:      name,
        DonorPhone:     phone,
        Amount:         amount,
        Purpose:        purpose,
        PaymentMode:    paymentApp, // Stores which app they used
        ReceiptNo:      receiptNo,
        Date:           time.Now(),
        Notes:          "paid", // Mark as paid since they are uploading screenshot
        ScreenshotPath: screenshotPath,
        TempleID:       1,
    }

    if err := database.GetDB().Create(&donation).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true, 
        "message": "Donation received successfully!",
        "receipt": receiptNo,
    })
}

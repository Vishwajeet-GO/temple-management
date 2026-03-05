package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GenerateReceipt(c *gin.Context) {
	id := c.Param("id")

	var donation models.Donation
	if database.DB.Preload("Festival").First(&donation, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Donation not found"})
		return
	}

	if donation.Status != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Donation not verified yet"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// ================= CERTIFICATE BORDER =================
	pdf.SetDrawColor(218, 165, 32) // Gold color
	pdf.SetLineWidth(2)
	pdf.Rect(8, 8, 194, 281, "D")

	pdf.SetLineWidth(0.5)
	pdf.Rect(12, 12, 186, 273, "D")

	// Reset line width for subsequent drawing
	pdf.SetLineWidth(0.2)

	// ================= HEADER BAND =================
	pdf.SetFillColor(255, 140, 0)
	pdf.Rect(12, 12, 186, 35, "F")

	// Logo - only load if the file exists (a failed Image call breaks the entire PDF)
	if _, err := os.Stat("./web/static/images/logo.png"); err == nil {
		pdf.Image("./web/static/images/logo.png", 18, 15, 20, 0, false, "", 0, "")
	}

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 22)
	pdf.SetXY(12, 20)
	pdf.CellFormat(186, 10, "Shree Gauri Shankar Mandir", "", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.SetXY(12, 30)
	pdf.CellFormat(186, 8, "Nirmal Singh Chawl, Gaondevi Road, Poisar, Kandivali East, Mumbai - Maharashtra", "", 0, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(60)

	// ================= TITLE =================
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "DONATION CERTIFICATE", "", 1, "C", false, 0, "")

	pdf.SetDrawColor(255, 140, 0)
	pdf.Line(70, pdf.GetY()+2, 140, pdf.GetY()+2)
	pdf.Ln(12)

	receiptNo := fmt.Sprintf("REC-%06d", donation.ID)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(100, 8, "Receipt No: "+receiptNo)
	pdf.Cell(0, 8, "Date: "+time.Now().Format("02 Jan 2006"))
	pdf.Ln(15)

	// ================= DETAILS =================
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 8, "Donation Details")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Donor Name:")
	pdf.Cell(0, 8, donation.Donor)
	pdf.Ln(8)

	pdf.Cell(60, 8, "Phone:")
	pdf.Cell(0, 8, donation.Phone)
	pdf.Ln(8)

	if donation.Festival != nil {
		pdf.Cell(60, 8, "Festival:")
		pdf.Cell(0, 8, donation.Festival.Name)
		pdf.Ln(8)
	}

	pdf.Cell(60, 8, "Payment Mode:")
	pdf.Cell(0, 8, donation.PaymentMode)
	pdf.Ln(12)

	// ================= AMOUNT BOX =================
	pdf.SetFillColor(255, 239, 213)
	pdf.Rect(50, pdf.GetY(), 110, 18, "F")

	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 128, 0)
	pdf.SetXY(50, pdf.GetY()+5)
	pdf.CellFormat(110, 10, fmt.Sprintf("Rs %.2f", donation.Amount), "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.Ln(30)

	// ================= THANK YOU =================
	pdf.SetFont("Arial", "I", 12)
	pdf.MultiCell(0, 8, "We sincerely thank you for your generous contribution towards the temple. May Lord Shiva bless you with prosperity and happiness.", "", "C", false)

	pdf.Ln(20)

	// ================= VERIFICATION =================
	pdf.SetFont("Arial", "", 9)
	verificationID := fmt.Sprintf("Verification ID: REC-%06d | Donation ID: %d", donation.ID, donation.ID)
	pdf.CellFormat(0, 6, verificationID, "", 1, "L", false, 0, "")

	// ================= SIGNATURE & STAMP =================
	pdf.SetXY(140, 220)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Authorized Signature")
	pdf.Line(135, 232, 190, 232)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Println("PDF generation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate receipt"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=Donation_Certificate.pdf")
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

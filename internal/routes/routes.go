package routes

import (
	"os"
	"strings"
	"time"

	"newapp/internal/handlers"
	"newapp/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://temple-management-o0yq.onrender.com"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	os.MkdirAll("./uploads/donations", 0755)
	os.MkdirAll("./uploads/expenses", 0755)

	r.Static("/static", "./web/static")
	r.Static("/uploads", "./uploads")

	// HTML pages
	r.GET("/", func(c *gin.Context) { c.File("./web/templates/index.html") })
	r.GET("/donations", func(c *gin.Context) { c.File("./web/templates/donations.html") })
	r.GET("/donate", func(c *gin.Context) { c.File("./web/templates/donate.html") })
	r.GET("/expenses", func(c *gin.Context) { c.File("./web/templates/expenses.html") })
	r.GET("/festivals", func(c *gin.Context) { c.File("./web/templates/festivals.html") })
	r.GET("/login", func(c *gin.Context) { c.File("./web/templates/login.html") })
	r.GET("/reports", func(c *gin.Context) { c.File("./web/templates/reports.html") })
	r.GET("/favicon.ico", func(c *gin.Context) { c.Status(204) })
	r.GET("/gallery", func(c *gin.Context) { c.File("./web/templates/gallery.html") })

	v1 := r.Group("/api/v1")
	{
		// Auth
		v1.POST("/login", handlers.Login)
		v1.POST("/logout", handlers.Logout)
		v1.GET("/auth/check", handlers.AuthCheck)

		// Dashboard
		v1.GET("/dashboard/summary", handlers.GetDashboardSummary)

		// Temple
		v1.GET("/temple", handlers.GetTemple)
		v1.PUT("/temple", middleware.AuthRequired("admin"), handlers.UpdateTemple)

		// Donations
		v1.GET("/donations", handlers.GetDonations)
		v1.POST("/donations", handlers.CreateDonation)
		v1.PUT("/donations/:id", middleware.AuthRequired("admin"), handlers.UpdateDonation)
		v1.DELETE("/donations/:id", middleware.AuthRequired("admin"), handlers.DeleteDonation)
		v1.PATCH("/donations/:id/toggle", middleware.AuthRequired("admin"), handlers.ToggleDonation)

		// Expenses
		v1.GET("/expenses", handlers.GetExpenses)
		v1.POST("/expenses", handlers.CreateExpense)
		v1.PUT("/expenses/:id", middleware.AuthRequired("admin"), handlers.UpdateExpense)
		v1.DELETE("/expenses/:id", middleware.AuthRequired("admin"), handlers.DeleteExpense)
		v1.PATCH("/expenses/:id/toggle", middleware.AuthRequired("admin"), handlers.ToggleExpense)

		// Festivals
		v1.GET("/festivals", handlers.GetFestivals)
		v1.GET("/festivals/:id/report", handlers.GetFestivalReport)
		v1.POST("/festivals", middleware.AuthRequired("admin"), handlers.CreateFestival)
		v1.PUT("/festivals/:id", middleware.AuthRequired("admin"), handlers.UpdateFestival)
		v1.DELETE("/festivals/:id", middleware.AuthRequired("admin"), handlers.DeleteFestival)

		// Gallery
		v1.GET("/gallery", handlers.GetGallery)
		v1.POST("/gallery", middleware.AuthRequired("admin"), handlers.UploadGallery)
		v1.DELETE("/gallery/:id", middleware.AuthRequired("admin"), handlers.DeleteGalleryItem)

		// Public donation (from donate page)
		v1.POST("/submit-donation", handlers.SubmitDonation)
		v1.GET("/payment-info", handlers.GetPaymentInfo)
	}

	// Admin
	v1.GET("/admin/dashboard", middleware.AuthRequired("admin"), handlers.GetAdminDashboard)
	v1.PUT("/admin/password", middleware.AuthRequired("admin"), handlers.UpdatePassword)

	r.GET("/admin", func(c *gin.Context) { c.File("./web/templates/admin.html") })
	// Backward compat
	r.GET("/api/donations", handlers.GetDonations)
	r.GET("/api/expenses", handlers.GetExpenses)
	r.GET("/api/festivals", handlers.GetFestivals)
	r.GET("/api/dashboard/summary", handlers.GetDashboardSummary)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{"success": false, "error": "Not found"})
			return
		}
		c.File("./web/templates/index.html")
	})

	return r
}

package routes

import (
	"net/http"
	"newapp/internal/handlers"
	"newapp/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/static", "./web/static")

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*.html")

	// Web routes - serve HTML files
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/festivals", func(c *gin.Context) {
		c.HTML(http.StatusOK, "festivals.html", nil)
	})

	r.GET("/donations", func(c *gin.Context) {
		c.HTML(http.StatusOK, "donations.html", nil)
	})

	r.GET("/expenses", func(c *gin.Context) {
		c.HTML(http.StatusOK, "expenses.html", nil)
	})

	r.GET("/reports", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reports.html", nil)
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth
		authHandler := handlers.NewAuthHandler()
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)

		// Dashboard
		dashboardHandler := handlers.NewDashboardHandler()
		api.GET("/dashboard/summary", dashboardHandler.GetSummary)
		api.GET("/dashboard/recent-donations", dashboardHandler.GetRecentDonations)
		api.GET("/dashboard/recent-expenses", dashboardHandler.GetRecentExpenses)

		// Temple
		templeHandler := handlers.NewTempleHandler()
		api.GET("/temple", templeHandler.GetTemple)

		// Festivals
		festivalHandler := handlers.NewFestivalHandler()
		api.GET("/festivals", festivalHandler.GetAll)
		api.GET("/festivals/upcoming", festivalHandler.GetUpcoming)
		api.GET("/festivals/:id", festivalHandler.GetByID)

		// Donations
		donationHandler := handlers.NewDonationHandler()
		api.GET("/donations", donationHandler.GetAll)
		api.GET("/donations/:id", donationHandler.GetByID)
		api.GET("/donations/festival/:festivalId", donationHandler.GetByFestival)

		// Expenses
		expenseHandler := handlers.NewExpenseHandler()
		api.GET("/expenses", expenseHandler.GetAll)
		api.GET("/expenses/:id", expenseHandler.GetByID)
		api.GET("/expenses/festival/:festivalId", expenseHandler.GetByFestival)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/auth/me", authHandler.GetCurrentUser)
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.PUT("/temple", templeHandler.UpdateTemple)
			protected.POST("/festivals", festivalHandler.Create)
			protected.PUT("/festivals/:id", festivalHandler.Update)
			protected.DELETE("/festivals/:id", festivalHandler.Delete)
			protected.POST("/donations", donationHandler.Create)
			protected.PUT("/donations/:id", donationHandler.Update)
			protected.DELETE("/donations/:id", donationHandler.Delete)
			protected.POST("/expenses", expenseHandler.Create)
			protected.PUT("/expenses/:id", expenseHandler.Update)
			protected.DELETE("/expenses/:id", expenseHandler.Delete)
		}
	}
}

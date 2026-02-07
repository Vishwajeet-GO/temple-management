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
	r.LoadHTMLGlob("web/templates/*.html")

	// Public Pages
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/festivals", func(c *gin.Context) { c.HTML(http.StatusOK, "festivals.html", nil) })
	r.GET("/donations", func(c *gin.Context) { c.HTML(http.StatusOK, "donations.html", nil) })
	r.GET("/expenses", func(c *gin.Context) { c.HTML(http.StatusOK, "expenses.html", nil) })
	r.GET("/reports", func(c *gin.Context) { c.HTML(http.StatusOK, "reports.html", nil) })

	// NEW: Donate Page
	r.GET("/donate", func(c *gin.Context) { c.HTML(http.StatusOK, "donate.html", nil) })

	api := r.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler()
		api.POST("/auth/login", authHandler.Login)

		// NEW: Payment Submit Route
		paymentHandler := handlers.NewPaymentHandler()
		api.POST("/submit-donation", paymentHandler.ProcessDonation)

		// Public Data
		templeHandler := handlers.NewTempleHandler()
		api.GET("/temple", templeHandler.GetTemple)

		festivalHandler := handlers.NewFestivalHandler()
		api.GET("/festivals", festivalHandler.GetAll)

		dashboardHandler := handlers.NewDashboardHandler()
		api.GET("/dashboard/summary", dashboardHandler.GetSummary)
		api.GET("/dashboard/recent-donations", dashboardHandler.GetRecentDonations)
		api.GET("/dashboard/recent-expenses", dashboardHandler.GetRecentExpenses)
		api.GET("/dashboard/project/:id", dashboardHandler.GetProjectStats)

		// Protected Routes
		donationHandler := handlers.NewDonationHandler()
		api.GET("/donations", donationHandler.GetAll) // Kept public for viewing based on previous requests

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Protected Management Routes
			api.POST("/donations", donationHandler.Create)
			api.PUT("/donations/:id", donationHandler.Update)
			api.DELETE("/donations/:id", donationHandler.Delete)

			expenseHandler := handlers.NewExpenseHandler()
			api.GET("/expenses", expenseHandler.GetAll)
			api.POST("/expenses", expenseHandler.Create)
			api.PUT("/expenses/:id", expenseHandler.Update)
			api.DELETE("/expenses/:id", expenseHandler.Delete)

			api.POST("/festivals", festivalHandler.Create)
			api.PUT("/festivals/:id", festivalHandler.Update)
			api.DELETE("/festivals/:id", festivalHandler.Delete)
		}
	}
}

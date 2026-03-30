package main

import (
	"log"

	"serverpoint/config"
	"serverpoint/database"
	"serverpoint/handlers"
	"serverpoint/middleware"
	"serverpoint/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)

	scheduler := services.NewScheduler(database.DB, cfg)
	scheduler.Start()
	defer scheduler.Stop()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authHandler := &handlers.AuthHandler{Config: cfg}
	serverHandler := &handlers.ServerHandler{Config: cfg}
	metricsHandler := &handlers.MetricsHandler{}

	api := r.Group("/api")
	{
		api.GET("/auth/check", authHandler.CheckSetup)
		api.POST("/auth/register", middleware.RateLimit(), authHandler.Register)
		api.POST("/auth/login", middleware.RateLimit(), authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthRequired(cfg))
		{
			protected.GET("/servers", serverHandler.List)
			protected.POST("/servers", serverHandler.Create)
			protected.GET("/servers/:id", serverHandler.Get)
			protected.DELETE("/servers/:id", serverHandler.Delete)

			protected.GET("/servers/:id/metrics", metricsHandler.GetMetrics)
			protected.GET("/servers/:id/metrics/latest", metricsHandler.GetLatestMetric)
			protected.GET("/servers/:id/processes", metricsHandler.GetProcesses)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

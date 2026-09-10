package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"caesar-quiz/internal/api"
	"caesar-quiz/internal/config"
	"caesar-quiz/internal/db"
	"caesar-quiz/internal/models"
	"caesar-quiz/internal/ws"
)

func main() {
	cfg := config.Load()

	// Init DB (pakai global instance)
	if err := db.InitDB(cfg.DatabaseDSN); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	database := db.GetDB()

	// Auto-migrate (development). For production use migrations SQL.
	if err := database.AutoMigrate(
		&models.Game{},
		&models.Player{},
		&models.Question{},
		&models.Answer{},
		&models.HallOfFame{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	hub := ws.NewHub(database)
	go hub.Run()

	r := gin.Default()

	// CORS dari env CORS_ORIGINS (dipisah koma). Default "*".
	allowCredentials := false
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			allowCredentials = false
			break
		}
		allowCredentials = true
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Host-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	api.RegisterRoutes(r, database, hub, cfg.AllowedOrigins)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s (allowed origins: %v)", port, cfg.AllowedOrigins)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

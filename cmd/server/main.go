package main
import (
	"log"

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

       // Enable CORS for frontend

       r.Use(cors.New(cors.Config{
	       AllowOrigins:     []string{"http://139.59.217.119:3000"},
	       AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	       AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
	       AllowCredentials: true,
       }))

       api.RegisterRoutes(r, database, hub)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

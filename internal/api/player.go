package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPlayerRoutes dipertahankan untuk kompatibilitas. Logika join sekarang
// terpusat di joinGameHandler (routes.go).
func RegisterPlayerRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/players/join", joinGameHandler(db))
}

package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB inisialisasi koneksi dan simpan di global DB
func InitDB(dsn string) error {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = database
	return nil
}

// GetDB return instance global *gorm.DB
func GetDB() *gorm.DB {
	return DB
}

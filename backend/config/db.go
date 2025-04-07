package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "myuser")
	password := getEnv("DB_PASSWORD", "mypassword")
	dbname := getEnv("DB_NAME", "mydatabase")
	port := getEnv("DB_PORT", "5432")
	timezone := getEnv("DB_TIMEZONE", "Asia/Shanghai")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		host, user, password, dbname, port, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic("failed to connect database")
	}

	log.Println("Connected to database successfully")

	db.AutoMigrate(&models.PatientResponse{})
	db.AutoMigrate(&models.Hospital{}, &models.Staff{}, &models.Patient{})
	db.AutoMigrate(&models.Token{})

	DB = db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

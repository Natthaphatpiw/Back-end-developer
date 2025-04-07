package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database with SQLite in-memory
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.Token{})
	db.AutoMigrate(&models.Staff{})
	db.AutoMigrate(&models.Hospital{})

	config.DB = db
	return db, nil
}

// SeedTestData inserts test data for middleware tests
func SeedTestData(db *gorm.DB) error {
	// Create test hospital
	hospital := models.Hospital{
		Name:     "Test Hospital",
		Location: "Test Location",
	}
	if err := db.Create(&hospital).Error; err != nil {
		return err
	}

	// Create test staff
	staff := models.Staff{
		Username:   "testuser",
		Password:   "hashedpassword",
		Name:       "Test User",
		Email:      "test@example.com",
		HospitalID: hospital.ID,
	}
	if err := db.Create(&staff).Error; err != nil {
		return err
	}

	// Create test tokens
	tokens := []models.Token{
		{
			Token:      "valid-token-12345",
			StaffID:    staff.ID,
			HospitalID: hospital.ID,
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		},
		{
			Token:      "expired-token-12345",
			StaffID:    staff.ID,
			HospitalID: hospital.ID,
			ExpiresAt:  time.Now().Add(-24 * time.Hour), // Expired
		},
	}

	for _, token := range tokens {
		if err := db.Create(&token).Error; err != nil {
			return err
		}
	}

	return nil
}

func TestAuthRequired(t *testing.T) {
	// Setup
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = SeedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	// Setup test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthRequired())

	router.GET("/protected", func(c *gin.Context) {
		staffID, exists := c.Get("staff_id")
		if !exists {
			c.JSON(500, gin.H{"error": "staff_id not set"})
			return
		}

		hospitalID, exists := c.Get("hospital_id")
		if !exists {
			c.JSON(500, gin.H{"error": "hospital_id not set"})
			return
		}

		c.JSON(200, gin.H{
			"message":     "success",
			"staff_id":    staffID,
			"hospital_id": hospitalID,
		})
	})

	// Test cases

	// Test case 1: No token provided
	t.Run("No Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Authentication required")
	})

	// Test case 2: Invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-999")
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid token")
	})

	// Test case 3: Expired token
	t.Run("Expired Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer expired-token-12345")
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Token expired")
	})

	// Test case 4: Valid token
	t.Run("Valid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token-12345")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["message"])
		assert.NotNil(t, response["staff_id"])
		assert.NotNil(t, response["hospital_id"])
	})
}

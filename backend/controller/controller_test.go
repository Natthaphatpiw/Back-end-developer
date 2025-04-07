package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/middleware"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
	db.AutoMigrate(&models.PatientResponse{})
	db.AutoMigrate(&models.Hospital{}, &models.Staff{}, &models.Patient{})
	db.AutoMigrate(&models.Token{})

	config.DB = db
	return db, nil
}

// SeedTestData inserts test data for unit tests
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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	staff := models.Staff{
		Username:   "testuser",
		Password:   string(hashedPassword),
		Name:       "Test User",
		Email:      "test@example.com",
		HospitalID: hospital.ID,
	}
	if err := db.Create(&staff).Error; err != nil {
		return err
	}

	// Create test token
	token := models.Token{
		Token:      "test-token-12345",
		StaffID:    staff.ID,
		HospitalID: hospital.ID,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(&token).Error; err != nil {
		return err
	}

	// Create expired token for testing
	expiredToken := models.Token{
		Token:      "expired-token-12345",
		StaffID:    staff.ID,
		HospitalID: hospital.ID,
		ExpiresAt:  time.Now().Add(-24 * time.Hour), // Expired
	}
	if err := db.Create(&expiredToken).Error; err != nil {
		return err
	}

	// Create test patients
	patients := []models.Patient{
		{
			FirstNameTh: "สมชาย",
			LastNameTh:  "ใจดี",
			FirstNameEn: "Somchai",
			LastNameEn:  "Jaidee",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			PatientHN:   "HN001",
			NationalID:  "1234567890123",
			PassportID:  "",
			PhoneNumber: "0891234567",
			Email:       "somchai@example.com",
			Gender:      "M",
			HospitalID:  hospital.ID,
		},
		{
			FirstNameTh: "สมหญิง",
			LastNameTh:  "รักดี",
			FirstNameEn: "Somying",
			LastNameEn:  "Rakdee",
			DateOfBirth: time.Date(1992, 5, 10, 0, 0, 0, 0, time.UTC),
			PatientHN:   "HN002",
			NationalID:  "1234567890124",
			PassportID:  "",
			PhoneNumber: "0891234568",
			Email:       "somying@example.com",
			Gender:      "F",
			HospitalID:  hospital.ID,
		},
	}

	for _, patient := range patients {
		if err := db.Create(&patient).Error; err != nil {
			return err
		}
	}

	return nil
}

// SetupRouter creates a test router with routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/patient/search/:id", GetPatient)

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/patient/search", SearchPatients)
	}

	router.POST("/staff/create", CreateStaff)
	router.POST("/staff/login", LoginStaff)

	return router
}

// TestGetPatient tests the GetPatient function
func TestGetPatient(t *testing.T) {
	// Setup
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = SeedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	router := SetupRouter()

	// Test case 1: Valid patient ID
	t.Run("Valid Patient ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search/1234567890123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response models.PatientResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "สมชาย", response.FirstNameTh)
		assert.Equal(t, "Jaidee", response.LastNameEn)
	})

	// Test case 2: Invalid patient ID
	t.Run("Invalid Patient ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search/9999999999999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// TestCreateStaff tests the CreateStaff function
func TestCreateStaff(t *testing.T) {
	// Setup
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = SeedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	router := SetupRouter()

	// Test case 1: Valid staff creation
	t.Run("Valid Staff Creation", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffCreateRequest{
			Username:   "newstaff",
			Password:   "password123",
			Name:       "New Staff",
			Email:      "newstaff@example.com",
			HospitalID: 1,
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		var response struct {
			Data models.StaffResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "newstaff", response.Data.Username)
	})

	// Test case 2: Duplicate username
	t.Run("Duplicate Username", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffCreateRequest{
			Username:   "testuser", // Already exists
			Password:   "password123",
			Name:       "Another Staff",
			Email:      "another@example.com",
			HospitalID: 1,
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Username already exists")
	})

	// Test case 3: Invalid hospital
	t.Run("Invalid Hospital", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffCreateRequest{
			Username:   "validstaff",
			Password:   "password123",
			Name:       "Valid Staff",
			Email:      "valid@example.com",
			HospitalID: 999, // Invalid hospital ID
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Hospital not found")
	})
}

// TestLoginStaff tests the LoginStaff function
func TestLoginStaff(t *testing.T) {
	// Setup
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = SeedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	router := SetupRouter()

	// Test case 1: Valid login
	t.Run("Valid Login", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffLoginRequest{
			Username:   "testuser",
			Password:   "password123",
			HospitalID: 1,
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response struct {
			Data models.TokenResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Data.Token)
		assert.Equal(t, "testuser", response.Data.Staff.Username)
	})

	// Test case 2: Invalid credentials
	t.Run("Invalid Password", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffLoginRequest{
			Username:   "testuser",
			Password:   "wrongpassword",
			HospitalID: 1,
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})

	// Test case 3: User not found
	t.Run("User Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()

		requestBody := models.StaffLoginRequest{
			Username:   "nonexistentuser",
			Password:   "password123",
			HospitalID: 1,
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})
}

// TestSearchPatients tests the SearchPatients function
func TestSearchPatients(t *testing.T) {
	// Setup
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = SeedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	router := SetupRouter()

	// Test case 1: Valid search with token
	t.Run("Valid Search", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search?first_name=สมชาย", nil)
		req.Header.Set("Authorization", "Bearer test-token-12345")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response struct {
			Data []models.PatientResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, "สมชาย", response.Data[0].FirstNameTh)
	})

	// Test case 2: No token provided
	t.Run("No Token Provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search?first_name=สมชาย", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Authentication required")
	})

	// Test case 3: Expired token
	t.Run("Expired Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search?first_name=สมชาย", nil)
		req.Header.Set("Authorization", "Bearer expired-token-12345")

		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Token expired")
	})

	// Test case 4: Multiple search parameters
	t.Run("Multiple Search Parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search?gender=F&phone_number=0891234568", nil)
		req.Header.Set("Authorization", "Bearer test-token-12345")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response struct {
			Data []models.PatientResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, "สมหญิง", response.Data[0].FirstNameTh)
	})

	// Test case 5: No results
	t.Run("No Results", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/patient/search?first_name=ไม่มีชื่อนี้", nil)
		req.Header.Set("Authorization", "Bearer test-token-12345")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response struct {
			Data []models.PatientResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Data))
	})
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

// Global TestDB variable
var TestDB *gorm.DB

// Test Check User Endpoint
func TestCheckUserHandler(t *testing.T) {
	app := &application{}

	t.Run("User Exists", func(t *testing.T) {
		// Ensure user exists
		DB.Where("email = ?", "existing@example.com").FirstOrCreate(&User{
			FullName: "Existing User", Email: "existing@example.com", Role: "Entrepreneur",
		})

		req, _ := http.NewRequest("GET", "/v1/user/check?email=existing@example.com", nil)
		rr := httptest.NewRecorder()

		app.checkUserHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"exists":true,"role":"Entrepreneur"}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("User Does Not Exist", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/user/check?email=newuser@example.com", nil)
		rr := httptest.NewRecorder()

		app.checkUserHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"exists":false}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Missing Email Parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/user/check", nil)
		rr := httptest.NewRecorder()

		app.checkUserHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expected := `{"status": "error", "message": "Email is required"}`
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

// Test Signup Handler
func TestSignupHandler(t *testing.T) {
	app := &application{}

	t.Run("Successful Signup", func(t *testing.T) {
		payload := SignupRequest{
			FullName:        "John Doe",
			Email:           "john.doe@example.com",
			Role:            "Entrepreneur",
			PhoneNumber:     "1234567890",
			Location:        "New York",
			LinkedInProfile: "https://linkedin.com/in/johndoe",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.signupHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"status":"success","message":"User registered successfully! Please verify your account."}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Duplicate Email Signup", func(t *testing.T) {
		// Ensure user already exists before attempting duplicate signup
		DB.Where("email = ?", "existing@example.com").FirstOrCreate(&User{
			FullName: "Existing User", Email: "existing@example.com", Role: "Entrepreneur",
		})

		payload := SignupRequest{
			FullName: "Existing User",
			Email:    "existing@example.com",
			Role:     "Entrepreneur",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.signupHandler(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code) // Expect conflict due to duplicate email
		expected := `{"status":"error","message":"User already exists with this email. Please log in."}`
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Invalid Request Payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.signupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expected := `{"status":"error","message":"Invalid request payload."}`
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

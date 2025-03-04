package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Initialize test database for connections tests
func setupTestDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // In-memory DB
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	err = DB.AutoMigrate(&User{}, &Entrepreneur{}, &Investor{}, &Mentor{}, &Startup{}, &Connection{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	fmt.Println("Test database initialized successfully!")
}

// Ensure database setup before running all tests in this file
func TestMain(m *testing.M) {
	setupTestDB() // Ensure test DB is ready before running tests
	code := m.Run()
	os.Exit(code)
}

func TestRequestConnectionHandler(t *testing.T) {
	app := &application{}

	t.Run("Successful Connection Request", func(t *testing.T) {
		// Ensure sender exists
		DB.Where("email = ?", "alice@example.com").FirstOrCreate(&User{
			FullName: "Alice", Email: "alice@example.com", Role: "Investor",
		})

		// Ensure receiver exists
		DB.Where("email = ?", "bob@example.com").FirstOrCreate(&User{
			FullName: "Bob", Email: "bob@example.com", Role: "Entrepreneur",
		})

		payload := ConnectionRequest{SenderEmail: "alice@example.com", ReceiverEmail: "bob@example.com"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/v1/connection/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.requestConnectionHandler(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), `"status":"success"`)
	})

	t.Run("Connection Already Exists", func(t *testing.T) {
		payload := ConnectionRequest{SenderEmail: "alice@example.com", ReceiverEmail: "bob@example.com"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/v1/connection/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.requestConnectionHandler(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("Invalid Sender Role", func(t *testing.T) {
		DB.Where("email = ?", "charlie@example.com").FirstOrCreate(&User{
			FullName: "Charlie", Email: "charlie@example.com", Role: "Entrepreneur",
		})

		payload := ConnectionRequest{SenderEmail: "charlie@example.com", ReceiverEmail: "bob@example.com"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/v1/connection/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.requestConnectionHandler(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestAcceptConnectionHandler(t *testing.T) {
	app := &application{}

	t.Run("Successful Connection Acceptance", func(t *testing.T) {
		var entrepreneur User
		DB.Where("email = ?", "bob@example.com").FirstOrCreate(&entrepreneur, User{FullName: "Bob", Email: "bob@example.com", Role: "Entrepreneur"})

		connection := Connection{SenderID: entrepreneur.ID + 1, ReceiverID: entrepreneur.ID, Status: "pending"}
		DB.Create(&connection)

		req, _ := http.NewRequest("PUT", "/v1/connection/accept/"+fmt.Sprint(connection.ID)+"?email=bob@example.com", nil)
		rr := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprint(connection.ID)}},
		}))

		app.acceptConnectionHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"status":"success"`)
	})

	t.Run("Reject Non-Existing Connection", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/connection/accept/999?email=bob@example.com", nil)
		rr := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{"999"}},
		}))

		app.acceptConnectionHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Unauthorized User Accepting Connection", func(t *testing.T) {
		DB.Where("email = ?", "eve@example.com").FirstOrCreate(&User{
			FullName: "Eve", Email: "eve@example.com", Role: "Investor",
		})

		var connection Connection
		DB.Where("status = ?", "pending").FirstOrCreate(&connection, Connection{SenderID: 1, ReceiverID: 2, Status: "pending"})

		req, _ := http.NewRequest("PUT", "/v1/connection/accept/"+fmt.Sprint(connection.ID)+"?email=eve@example.com", nil)
		rr := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprint(connection.ID)}},
		}))

		app.acceptConnectionHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestRejectConnectionHandler(t *testing.T) {
	app := &application{}

	t.Run("Successful Connection Rejection", func(t *testing.T) {
		var entrepreneur User
		DB.Where("email = ?", "bob@example.com").FirstOrCreate(&entrepreneur, User{
			FullName: "Bob", Email: "bob@example.com", Role: "Entrepreneur",
		})

		// Ensure sender exists
		var sender User
		DB.Where("email = ?", "alice@example.com").FirstOrCreate(&sender, User{
			FullName: "Alice", Email: "alice@example.com", Role: "Investor",
		})

		// Create a pending connection request
		connection := Connection{SenderID: sender.ID, ReceiverID: entrepreneur.ID, Status: "pending"}
		DB.Create(&connection)

		// Create request to reject connection
		req, _ := http.NewRequest("DELETE", "/v1/connection/reject/"+fmt.Sprint(connection.ID)+"?email=bob@example.com", nil)
		rr := httptest.NewRecorder()

		// Inject chi router param
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprint(connection.ID)}},
		}))

		app.rejectConnectionHandler(rr, req)

		// Ensure response is HTTP 200 (OK)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Unmarshal JSON response correctly (single object, not a list)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "success", responseBody["status"])
		assert.Contains(t, responseBody["message"], "Connection request")
	})

	t.Run("Reject Non-Existing Connection", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/connection/reject/999?email=bob@example.com", nil)
		rr := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{"999"}},
		}))

		app.rejectConnectionHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestGetPendingConnectionsHandler(t *testing.T) {
	app := &application{}

	t.Run("Retrieve Pending Connections Successfully", func(t *testing.T) {
		// Ensure entrepreneur exists
		var entrepreneur User
		DB.Where("email = ?", "bob@example.com").FirstOrCreate(&entrepreneur, User{
			FullName: "Bob", Email: "bob@example.com", Role: "Entrepreneur",
		})

		// Ensure sender (Investor) exists
		var sender User
		DB.Where("email = ?", "alice@example.com").FirstOrCreate(&sender, User{
			FullName: "Alice", Email: "alice@example.com", Role: "Investor",
		})

		// Create a pending connection
		connection := Connection{SenderID: sender.ID, ReceiverID: entrepreneur.ID, Status: "pending"}
		DB.Create(&connection)

		// Make API request
		req, _ := http.NewRequest("GET", "/v1/connection/pending?email=bob@example.com", nil)
		rr := httptest.NewRecorder()

		app.getPendingConnectionsHandler(rr, req)

		// Ensure HTTP 200 OK response
		assert.Equal(t, http.StatusOK, rr.Code)

		// Unmarshal response correctly
		var responseBody []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)

		// Ensure at least one pending connection exists
		assert.Greater(t, len(responseBody), 0)

		// Check if it contains expected user details
		assert.Equal(t, "alice@example.com", responseBody[0]["email"])
		assert.Equal(t, "Alice", responseBody[0]["full_name"])
		assert.Equal(t, "Investor", responseBody[0]["role"])
	})
}

func TestGetAcceptedConnectionsHandler(t *testing.T) {
	app := &application{}

	t.Run("Retrieve Accepted Connections Successfully", func(t *testing.T) {
		DB.Where("email = ?", "alice@example.com").FirstOrCreate(&User{
			FullName: "Alice", Email: "alice@example.com", Role: "Investor",
		})

		req, _ := http.NewRequest("GET", "/v1/connection/my?email=alice@example.com", nil)
		rr := httptest.NewRecorder()

		app.getAcceptedConnectionsHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"status":"success"`)
	})
}

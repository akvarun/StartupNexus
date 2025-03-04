package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ConnectionRequest struct {
	SenderEmail   string `json:"sender_email"`
	ReceiverEmail string `json:"receiver_email"`
}

// API handler to request a connection
func (app *application) requestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var req ConnectionRequest

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find sender (Investor or Mentor)
	var sender User
	if err := DB.Where("email = ?", req.SenderEmail).First(&sender).Error; err != nil {
		http.Error(w, "Sender not found", http.StatusNotFound)
		return
	}

	// Find receiver (Entrepreneur)
	var receiver User
	if err := DB.Where("email = ?", req.ReceiverEmail).First(&receiver).Error; err != nil {
		http.Error(w, "Receiver not found", http.StatusNotFound)
		return
	}

	// Ensure sender is Investor/Mentor and receiver is Entrepreneur
	if sender.Role != "Investor" && sender.Role != "Mentor" {
		http.Error(w, "Only Investors and Mentors can send connection requests", http.StatusForbidden)
		return
	}
	if receiver.Role != "Entrepreneur" {
		http.Error(w, "Only Entrepreneurs can receive connection requests", http.StatusForbidden)
		return
	}

	// Check if a connection already exists
	var existingConnection Connection
	if err := DB.Where("sender_id = ? AND receiver_id = ?", sender.ID, receiver.ID).First(&existingConnection).Error; err == nil {
		http.Error(w, "Connection already exists", http.StatusConflict)
		return
	}

	// Create a new connection request
	newConnection := Connection{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Status:     "pending",
	}
	if err := DB.Create(&newConnection).Error; err != nil {
		http.Error(w, "Failed to create connection request", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Connection request sent from %s to %s", sender.Email, receiver.Email),
	})
}

func (app *application) acceptConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Get connection ID from URL params
	connID := chi.URLParam(r, "id")

	// Get entrepreneur email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find the entrepreneur by email
	var entrepreneur User
	if err := DB.Where("email = ? AND role = ?", email, "Entrepreneur").First(&entrepreneur).Error; err != nil {
		http.Error(w, "Entrepreneur not found", http.StatusNotFound)
		return
	}

	// Find the pending connection request
	var connection Connection
	err := DB.Where("id = ? AND receiver_id = ? AND status = ?", connID, entrepreneur.ID, "pending").First(&connection).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Pending connection request not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Update connection status to accepted
	connection.Status = "accepted"
	if err := DB.Save(&connection).Error; err != nil {
		http.Error(w, "Failed to accept connection", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Connection request %s accepted", connID),
	})
}

func (app *application) getPendingConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get entrepreneur email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find the entrepreneur by email
	var entrepreneur User
	if err := DB.Where("email = ? AND role = ?", email, "Entrepreneur").First(&entrepreneur).Error; err != nil {
		http.Error(w, "Entrepreneur not found", http.StatusNotFound)
		return
	}

	// Fetch pending connection requests
	var pendingConnections []Connection
	if err := DB.Where("receiver_id = ? AND status = ?", entrepreneur.ID, "pending").Find(&pendingConnections).Error; err != nil {
		http.Error(w, "Error fetching pending connections", http.StatusInternalServerError)
		return
	}

	// Retrieve details of senders (Investors/Mentors)
	var connectionDetails []map[string]interface{}
	for _, conn := range pendingConnections {
		var sender User
		if err := DB.Where("id = ?", conn.SenderID).First(&sender).Error; err == nil {
			connectionDetails = append(connectionDetails, map[string]interface{}{
				"connection_id": conn.ID,
				"full_name":     sender.FullName,
				"email":         sender.Email,
				"role":          sender.Role,
				"location":      sender.Location,
				"linkedin":      sender.LinkedIn,
			})
		}
	}

	// Return pending connections
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connectionDetails)
}

func (app *application) rejectConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Get connection ID from URL params
	connID := chi.URLParam(r, "id")

	// Get entrepreneur email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find the entrepreneur by email
	var entrepreneur User
	if err := DB.Where("email = ? AND role = ?", email, "Entrepreneur").First(&entrepreneur).Error; err != nil {
		http.Error(w, "Entrepreneur not found", http.StatusNotFound)
		return
	}

	// Find the pending connection request
	var connection Connection
	err := DB.Where("id = ? AND receiver_id = ? AND status = ?", connID, entrepreneur.ID, "pending").First(&connection).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Pending connection request not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	DB.Model(&connection).Update("status", "rejected")

	// Delete the connection request
	if err := DB.Delete(&connection).Error; err != nil {
		http.Error(w, "Failed to reject connection request", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Connection request %s rejected", connID),
	})
}

func (app *application) getAcceptedConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch all accepted connections for the user
	var connections []Connection
	if err := DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", user.ID, user.ID, "accepted").Find(&connections).Error; err != nil {
		http.Error(w, "Failed to retrieve connections", http.StatusInternalServerError)
		return
	}

	// Create response structure
	var result []map[string]interface{}
	for _, conn := range connections {
		var connectedUser User
		if conn.SenderID == user.ID {
			// If user is the sender, get receiver details
			DB.Where("id = ?", conn.ReceiverID).First(&connectedUser)
		} else {
			// If user is the receiver, get sender details
			DB.Where("id = ?", conn.SenderID).First(&connectedUser)
		}

		// Append connection details
		result = append(result, map[string]interface{}{
			"id":    connectedUser.ID,
			"name":  connectedUser.FullName,
			"email": connectedUser.Email,
			"role":  connectedUser.Role,
		})
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"connections": result,
	})
}

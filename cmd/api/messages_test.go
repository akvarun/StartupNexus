package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleSendMessage(t *testing.T) {
	setupTestDB()
	app := &application{}

	t.Run("Send Message Successfully", func(t *testing.T) {
		// Create users
		sender := User{FullName: "Alice", Email: "alice@example.com", Role: "Entrepreneur"}
		receiver := User{FullName: "Bob", Email: "bob@example.com", Role: "Investor"}
		DB.Create(&sender)
		DB.Create(&receiver)

		// Create accepted connection
		connection := Connection{SenderID: sender.ID, ReceiverID: receiver.ID, Status: "accepted"}
		DB.Create(&connection)

		// Prepare payload
		payload := map[string]string{
			"sender_email":   "alice@example.com",
			"receiver_email": "bob@example.com",
			"content":        "Hello, Bob!",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/message/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.handleSendMessage(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Message sent successfully")
	})

	t.Run("Send Fails - Not Connected", func(t *testing.T) {
		payload := map[string]string{
			"sender_email":   "stranger@example.com",
			"receiver_email": "bob@example.com",
			"content":        "Hi Bob!",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/message/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.handleSendMessage(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Sender not found")
	})
}

func TestHandleGetMessages(t *testing.T) {
	setupTestDB()
	app := &application{}

	t.Run("Get Chat Messages", func(t *testing.T) {
		// Create users
		user1 := User{FullName: "Charlie", Email: "charlie@example.com", Role: "Mentor"}
		user2 := User{FullName: "Diana", Email: "diana@example.com", Role: "Entrepreneur"}
		DB.Create(&user1)
		DB.Create(&user2)

		// Create accepted connection
		connection := Connection{SenderID: user1.ID, ReceiverID: user2.ID, Status: "accepted"}
		DB.Create(&connection)

		// Add messages
		DB.Create(&Message{SenderEmail: user1.Email, ReceiverEmail: user2.Email, Content: "Hi Diana!"})
		DB.Create(&Message{SenderEmail: user2.Email, ReceiverEmail: user1.Email, Content: "Hello Charlie!"})

		req, _ := http.NewRequest("GET", "/v1/message/chat?sender=charlie@example.com&receiver=diana@example.com", nil)
		rr := httptest.NewRecorder()

		app.handleGetMessages(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Hi Diana!")
		assert.Contains(t, rr.Body.String(), "Hello Charlie!")
	})
}

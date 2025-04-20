package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

// handleSendMessage handles sending a message between connected users
func (app *application) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		SenderEmail   string `json:"sender_email"`
		ReceiverEmail string `json:"receiver_email"`
		Content       string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Validate users
	var sender, receiver User
	if err := DB.Where("email = ?", req.SenderEmail).First(&sender).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Sender not found"}`, http.StatusNotFound)
		return
	}
	if err := DB.Where("email = ?", req.ReceiverEmail).First(&receiver).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Receiver not found"}`, http.StatusNotFound)
		return
	}

	// Ensure they are connected
	var connection Connection
	if err := DB.Where(
		"((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?",
		sender.ID, receiver.ID, receiver.ID, sender.ID, "accepted").First(&connection).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Users are not connected"}`, http.StatusForbidden)
		return
	}

	// Store the message
	message := Message{
		SenderEmail:   req.SenderEmail,
		ReceiverEmail: req.ReceiverEmail,
		Content:       req.Content,
	}
	if err := DB.Create(&message).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Failed to save message"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Message sent successfully",
	})
}

// handleGetMessages retrieves chat history between two connected users
func (app *application) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	senderEmail := r.URL.Query().Get("sender")
	receiverEmail := r.URL.Query().Get("receiver")
	if senderEmail == "" || receiverEmail == "" {
		http.Error(w, `{"status":"error","message":"Missing sender or receiver email"}`, http.StatusBadRequest)
		return
	}

	var sender, receiver User
	if err := DB.Where("email = ?", senderEmail).First(&sender).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Sender not found"}`, http.StatusNotFound)
		return
	}
	if err := DB.Where("email = ?", receiverEmail).First(&receiver).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Receiver not found"}`, http.StatusNotFound)
		return
	}

	// Ensure they are connected
	var connection Connection
	if err := DB.Where(
		"((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?",
		sender.ID, receiver.ID, receiver.ID, sender.ID, "accepted").First(&connection).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Users are not connected"}`, http.StatusForbidden)
		return
	}

	// Retrieve messages in both directions
	var messages []Message
	if err := DB.Where(
		"(sender_email = ? AND receiver_email = ?) OR (sender_email = ? AND receiver_email = ?)",
		senderEmail, receiverEmail, receiverEmail, senderEmail).Find(&messages).Error; err != nil {
		http.Error(w, `{"status":"error","message":"Failed to retrieve messages"}`, http.StatusInternalServerError)
		return
	}

	// Sort messages by timestamp ascending
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	chat := []map[string]interface{}{}
	for _, m := range messages {
		chat = append(chat, map[string]interface{}{
			"from":      m.SenderEmail,
			"to":        m.ReceiverEmail,
			"content":   m.Content,
			"timestamp": m.CreatedAt,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"chat_with": receiverEmail,
		"messages":  chat,
	})
}

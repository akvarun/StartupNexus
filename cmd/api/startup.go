package main

import (
	"encoding/json"
	"net/http"
)

// Handler to fetch all startups
func (app *application) getAllStartupsHandler(w http.ResponseWriter, r *http.Request) {
	var startups []Startup
	if err := DB.Find(&startups).Error; err != nil {
		http.Error(w, "Failed to retrieve startups", http.StatusInternalServerError)
		return
	}

	var result []map[string]interface{}
	for _, startup := range startups {
		var entrepreneur User
		if err := DB.Where("id = ?", startup.EntrepreneurID).First(&entrepreneur).Error; err != nil {
			continue // Skip if entrepreneur not found
		}

		result = append(result, map[string]interface{}{
			"id":           startup.ID,
			"startup_name": startup.StartupName,
			"industry":     startup.Industry,
			"description":  startup.Description,
			"budget":       startup.Budget,
			"timeframe":    startup.Timeframe,
			"entrepreneur": map[string]interface{}{
				"id":       entrepreneur.ID,
				"name":     entrepreneur.FullName,
				"email":    entrepreneur.Email,
				"location": entrepreneur.Location,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"startups": result,
	})
}

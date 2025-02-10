package main

import (
	"encoding/json"
	"net/http"
)

// Request structure
type SignupRequest struct {
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	PhoneNumber       string `json:"phone_number"`
	Location          string `json:"location"`
	LinkedInProfile   string `json:"linkedin_profile"`
	VerificationProof string `json:"verification_proof,omitempty"`

	// Role-specific fields
	StartupName         string  `json:"startup_name,omitempty"`
	Industry            string  `json:"industry,omitempty"`
	Description         string  `json:"description,omitempty"`
	Budget              float64 `json:"budget,omitempty"`
	Timeframe           string  `json:"timeframe,omitempty"`
	InvestmentRange     string  `json:"investment_range,omitempty"`
	PreferredIndustries string  `json:"preferred_industries,omitempty"`
	ExperienceYears     int     `json:"experience_years,omitempty"`
	Expertise           string  `json:"expertise,omitempty"`
	PastMentorships     string  `json:"past_mentorships,omitempty"`
	YearsOfExperience   int     `json:"years_of_experience,omitempty"`
}

// Signup API
func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid request payload.",
		})
		return
	}
	// Check if the email already exists
	var existingUser User
	if err := DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "User already exists with this email. Please log in.",
		})
		return
	}

	user := User{
		FullName:           req.FullName,
		Email:              req.Email,
		Role:               req.Role,
		PhoneNumber:        req.PhoneNumber,
		Location:           req.Location,
		LinkedIn:           req.LinkedInProfile,
		VerificationProof:  &req.VerificationProof,
		VerificationStatus: "pending",
	}

	result := DB.Create(&user)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Error creating user.",
		})
		return
	}

	switch req.Role {
	case "Entrepreneur":
		entrepreneur := Entrepreneur{
			UserID:      user.ID,
			StartupName: req.StartupName,
			Industry:    req.Industry,
			Description: req.Description,
			Budget:      req.Budget,
			Timeframe:   req.Timeframe,
		}
		DB.Create(&entrepreneur)

	case "Investor":
		investor := Investor{
			UserID:              user.ID,
			InvestmentRange:     req.InvestmentRange,
			PreferredIndustries: req.PreferredIndustries,
			ExperienceYears:     req.ExperienceYears,
			VerificationProof:   req.VerificationProof,
		}
		DB.Create(&investor)

	case "Mentor":
		mentor := Mentor{
			UserID:            user.ID,
			Expertise:         req.Expertise,
			PastMentorships:   req.PastMentorships,
			YearsOfExperience: req.YearsOfExperience,
			VerificationProof: req.VerificationProof,
		}
		DB.Create(&mentor)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User registered successfully! Please verify your account.",
	})
}

// getting the users from the database
func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User

	// Fetch all users from the database
	if err := DB.Find(&users).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Failed to retrieve users.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

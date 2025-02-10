package main

import "gorm.io/gorm"

// Users Table
type User struct {
	gorm.Model
	FullName           string  `json:"full_name"`
	Email              string  `json:"email" gorm:"unique"`
	Role               string  `json:"role"`
	PhoneNumber        string  `json:"phone_number"`
	Location           string  `json:"location"`
	LinkedIn           string  `json:"linkedin_profile"`
	VerificationProof  *string `json:"verification_proof,omitempty"`
	VerificationStatus string  `json:"verification_status" gorm:"default:'pending'"`
}

// Entrepreneurs Table
type Entrepreneur struct {
	gorm.Model
	UserID      uint    `json:"user_id"`
	StartupName string  `json:"startup_name"`
	Industry    string  `json:"industry"`
	Description string  `json:"description"`
	Budget      float64 `json:"budget"`
	Timeframe   string  `json:"timeframe"`
}

// Investors Table
type Investor struct {
	gorm.Model
	UserID              uint   `json:"user_id"`
	InvestmentRange     string `json:"investment_range"`
	PreferredIndustries string `json:"preferred_industries"`
	ExperienceYears     int    `json:"experience_years"`
	VerificationProof   string `json:"verification_proof"`
}

// Mentors Table
type Mentor struct {
	gorm.Model
	UserID            uint   `json:"user_id"`
	Expertise         string `json:"expertise"`
	PastMentorships   string `json:"past_mentorships"`
	YearsOfExperience int    `json:"years_of_experience"`
	VerificationProof string `json:"verification_proof"`
}

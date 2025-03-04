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
	UserID uint `json:"user_id" gorm:"not null;unique;foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// Investors Table
type Investor struct {
	gorm.Model
	UserID              uint   `json:"user_id" gorm:"not null;unique;foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	InvestmentRange     string `json:"investment_range"`
	PreferredIndustries string `json:"preferred_industries"`
	ExperienceYears     int    `json:"experience_years"`
	VerificationProof   string `json:"verification_proof"`
}

// Mentors Table
type Mentor struct {
	gorm.Model
	UserID            uint   `json:"user_id" gorm:"not null;unique;foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Expertise         string `json:"expertise"`
	PastMentorships   string `json:"past_mentorships"`
	YearsOfExperience int    `json:"years_of_experience"`
	VerificationProof string `json:"verification_proof"`
}

// Startups Table
type Startup struct {
	gorm.Model
	EntrepreneurID uint    `json:"entrepreneur_id" gorm:"not null;foreignKey:EntrepreneurID;constraint:OnDelete:CASCADE;"`
	StartupName    string  `json:"startup_name" gorm:"not null"`
	Industry       string  `json:"industry"`
	Description    string  `json:"description"`
	Budget         float64 `json:"budget"`
	Timeframe      string  `json:"timeframe"`
}

// Connections Table (New)
type Connection struct {
	gorm.Model
	SenderID   uint   `json:"sender_id" gorm:"not null"`       // Investor or Mentor
	ReceiverID uint   `json:"receiver_id" gorm:"not null"`     // Entrepreneur
	Status     string `json:"status" gorm:"default:'pending'"` // pending, accepted, rejected
}

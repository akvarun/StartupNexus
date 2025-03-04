package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Global variable for database connection
var DB *gorm.DB

// Initialize the database connection
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate tables
	err = DB.AutoMigrate(&User{}, &Entrepreneur{}, &Investor{}, &Mentor{}, &Startup{}, &Connection{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration successful!")
}

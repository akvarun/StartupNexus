package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	err = DB.AutoMigrate(
		&User{},
		&Entrepreneur{},
		&Investor{},
		&Mentor{},
		&Startup{},
		&Connection{},
		&Comment{},
		&Like{},
		&Report{},
		&Message{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
}

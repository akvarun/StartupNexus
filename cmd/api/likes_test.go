package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"context"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForLikes() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	DB.AutoMigrate(&User{}, &Entrepreneur{}, &Startup{}, &Like{})
	fmt.Println("Test DB setup for likes completed")
}

func createTestUser(email, role string) User {
	user := User{
		FullName: fmt.Sprintf("Test %s", role),
		Email:    email,
		Role:     role,
	}
	DB.Create(&user)
	return user
}

func createTestStartup(user User) Startup {
	entrepreneur := Entrepreneur{UserID: user.ID}
	DB.Create(&entrepreneur)

	startup := Startup{
		EntrepreneurEmail: user.Email,
		EntrepreneurID:    entrepreneur.ID,
		StartupName:       "Test Startup",
		Industry:          "Tech",
		Description:       "A test startup",
		Budget:            10000,
		Timeframe:         "6 months",
	}
	DB.Create(&startup)
	return startup
}

func TestLikeHandlers(t *testing.T) {
	setupTestDBForLikes()
	app := &application{}

	t.Run("Like Startup Successfully", func(t *testing.T) {
		user := createTestUser("likeuser@example.com", "Entrepreneur")
		startup := createTestStartup(user)

		payload := map[string]interface{}{
			"email":      user.Email,
			"startup_id": startup.ID,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/startup/like", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.likeStartupHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Startup liked successfully")
	})

	t.Run("Get Startup Like Count", func(t *testing.T) {
		user := createTestUser("countuser@example.com", "Entrepreneur")
		startup := createTestStartup(user)
		DB.Create(&Like{UserID: user.ID, StartupID: startup.ID})

		req, _ := http.NewRequest("GET", "/v1/startup/likes/"+strconv.Itoa(int(startup.ID)), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", strconv.Itoa(int(startup.ID)))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		rr := httptest.NewRecorder()
		app.getStartupLikesHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "like_count")
	})

	t.Run("Unlike Startup Successfully", func(t *testing.T) {
		user := createTestUser("unlikeuser@example.com", "Entrepreneur")
		startup := createTestStartup(user)
		DB.Create(&Like{UserID: user.ID, StartupID: startup.ID})

		payload := map[string]interface{}{
			"email":      user.Email,
			"startup_id": startup.ID,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("DELETE", "/v1/startup/unlike", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		app.unlikeStartupHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Startup unliked successfully")
	})

	// t.Run("Like Startup - Already Liked", func(t *testing.T) {
	// 	user := createTestUser("dupelike@example.com", "Entrepreneur")
	// 	startup := createTestStartup(user)
	// 	DB.Create(&Like{UserID: user.ID, StartupID: startup.ID})

	// 	payload := map[string]interface{}{
	// 		"email":      user.Email,
	// 		"startup_id": startup.ID,
	// 	}
	// 	body, _ := json.Marshal(payload)

	// 	req, _ := http.NewRequest("POST", "/v1/startup/like", bytes.NewBuffer(body))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	rr := httptest.NewRecorder()

	// 	app.likeStartupHandler(rr, req)

	// 	assert.Equal(t, http.StatusConflict, rr.Code)
	// 	assert.Contains(t, rr.Body.String(), "Already liked")
	// })
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type config struct {
	addr string
}

type application struct {
	config config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/signup", app.signupHandler)
		r.Get("/user/check", app.checkUserHandler)
		r.Get("/user/search", app.searchUsersByNameHandler)
		r.Get("/connection/pending", app.getPendingConnectionsHandler)
		r.Get("/suggestions", app.getConnectionSuggestionsHandler)
		r.Post("/connection/request", app.requestConnectionHandler)
		r.Put("/connection/accept/{id}", app.acceptConnectionHandler)
		r.Delete("/connection/reject/{id}", app.rejectConnectionHandler)
		r.Get("/connection/my", app.getAcceptedConnectionsHandler)
		r.Get("/startup/all", app.getAllStartupsHandler)
		r.Get(("/startup/user"), app.getStartupsByUserHandler)
		r.Post("/startup/insert", app.insertStartupHandler)
		r.Put("/startup/update", app.updateStartupHandler)
		r.Delete("/startup/delete", app.deleteStartupHandler)
		r.Post("/startup/like", app.likeStartupHandler)
		r.Delete("/startup/unlike", app.unlikeStartupHandler)
		r.Get("/startup/likes/{id}", app.getStartupLikesHandler)
		r.Post("/startup/comment", app.addCommentHandler)
		r.Get("/startup/comments/{id}", app.getCommentsHandler)
		r.Post("/startup/report", app.reportStartupHandler)
		r.Post("/message/send", app.handleSendMessage)
		r.Get("/message/inbox", app.handleGetMessages)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Starting server on %s", srv.Addr)

	return srv.ListenAndServe()
}

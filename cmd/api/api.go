package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r.Route("/v1", func(r chi.Router) {
		r.Post("/signup", app.signupHandler)
		r.Get("/user/check", app.checkUserHandler)
		r.Get("/connection/pending", app.getPendingConnectionsHandler)
		r.Post("/connection/request", app.requestConnectionHandler)
		r.Put("/connection/accept/{id}", app.acceptConnectionHandler)
		r.Delete("/connection/reject/{id}", app.rejectConnectionHandler)
		r.Get("/connection/my", app.getAcceptedConnectionsHandler)
		r.Get("/startup/all", app.getAllStartupsHandler)
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

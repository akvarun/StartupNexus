package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	return mux
}

func (app *application) run(mux *http.ServeMux) error {
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

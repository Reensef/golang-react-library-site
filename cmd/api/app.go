package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Reensef/golang-react-boolib/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config appConfig
	store  store.Storage
}

type appConfig struct {
	appAddr      string
	sqlDBConfig  sqlDBConfig
	blobDBConfig blobDBConfig
	env          string
}

// TODO : move to db package
type sqlDBConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type blobDBConfig struct {
	addr string
	id   string
	key  string
}

var (
	ErrInvalidCommentID = errors.New("invalid comment ID")
	ErrInvalidPostID    = errors.New("invalid post ID")
	ErrInvalidLimit     = errors.New("invalid limit")
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.appAddr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started as %s", app.config.appAddr)

	return srv.ListenAndServe()
}
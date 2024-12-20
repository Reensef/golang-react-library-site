package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	host         string
	port         string
	name         string
	user         string
	password     string
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
	ErrInvalidFileID = errors.New("invalid file ID")
	ErrInvalidTagID  = errors.New("invalid tag ID")
)

type contextKey string

const (
	userInfoCtxKey contextKey = "userInfo"
)

type UserCtxInfo struct {
	ID   int64
	Role string
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обработка preflight-запросов OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// r.Use(app.enableCORS)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/login", app.loginHandler)
	r.Handle("/static/*", http.FileServer(http.Dir(app.static())))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, app.static()+"/index.html")
	})

	// TODO Проверить возвраты состояний
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(app.JWTMiddleware)
		r.Get("/health", app.healthCheckHandler)

		r.Route("/files", func(r chi.Router) {
			r.Post("/upload", app.uploadFileHandler)

			r.Get("/", app.getFilesHandler)
			r.Get("/{id}", app.getFileHandler)
			r.Get("/download/{id}", app.downloadFileHandler)
			r.Get("/open/{id}", app.openFileHandler)

			r.Get("/tags", app.getTagsHandler)

			r.Delete("/{id}", app.deleteFileHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Get("/{id}", app.getUserHandler)
		})
	})

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

func (app *application) static() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to find executable path: %v", err)
	}
	return filepath.Join(filepath.Dir(ex), "../client/dist")
}

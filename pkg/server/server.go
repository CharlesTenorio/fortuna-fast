package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/katana/fortuna/backend-go/internal/config"
)

func NewHTTPServer(r chi.Router, conf *config.Config) *http.Server {
	// Add middleware for logging, timeouts, etc.

	srv := &http.Server{
		ReadTimeout:  10 * time.Second, // Wait for 10 seconds for a request to be fully read
		WriteTimeout: 10 * time.Second, // Respond within 10 seconds
		Addr:         ":" + conf.PORT,
		Handler:      r,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}

	return srv
}

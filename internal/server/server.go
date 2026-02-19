package server

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/go-chi/chi/v5"

	"github.com/caddyserver/certmagic"
	"github.com/circa10a/go-rest-template/internal/server/handlers"
	"github.com/circa10a/go-rest-template/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed api.html
var apiDocs []byte

// Server is our web server that runs the network mirror.
type Server struct {
	// Embed configuration used to build the server.
	Config

	mux         http.Handler
	logger      *slog.Logger
	middlewares []func(http.Handler) http.Handler
}

// Config holds configuration for creating a Server.
type Config struct {
	TLSCert    string
	TLSKey     string
	LogFormat  string
	LogLevel   string
	Domains    []string
	Port       int
	AutoTLS    bool
	Metrics    bool
	Validation bool
}

// New returns a new server configured from cfg.
func New(cfg *Config) (*Server, error) {
	server := &Server{
		Config: *cfg,
	}

	if server.LogLevel == "" {
		server.LogLevel = "info"
	}

	server.LogFormat = strings.ToLower(server.LogFormat)

	router := chi.NewRouter()
	server.mux = router

	// Ensure configuration options are valid/compatible
	err := server.validate()
	if err != nil {
		return nil, err
	}

	logLevel, err := log.ParseLevel(server.LogLevel)
	if err != nil {
		return nil, err
	}

	logHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Formatter:       getLogFormatter(server.LogFormat),
		Level:           logLevel,
	})
	server.logger = slog.New(logHandler)

	// Features
	if server.Metrics {
		router.Handle("/metrics", promhttp.Handler())
		server.middlewares = append(server.middlewares, middleware.Prometheus)
	}

	// Default middlewares
	server.mux = middleware.Logging(server.logger, server.mux)

	// Add middlewares via http.Handler chaining
	for _, mw := range server.middlewares {
		server.mux = mw(server.mux)
	}

	// Routes
	router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(apiDocs) })
	router.Get("/health", handlers.HealthHandleFunc)

	return server, nil
}

// Start starts the listener of the server.
func (s *Server) Start() error {
	log := s.logger.With("component", "server")

	// Auto TLS will create listeners on port 80 and 443
	if s.AutoTLS {
		log.Info("Starting server on :80 and :443")
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = "user@oss.com"
		return certmagic.HTTPS(s.Domains, s.mux)
	}

	// If no auto TLS, use specified server port
	// :{port}
	addr := fmt.Sprintf(":%d", s.Port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	log.Info("Starting server on " + addr)

	// If custom cert and key provided, listen on specified server port via https
	if s.TLSCert != "" && s.TLSKey != "" {
		return httpServer.ListenAndServeTLS(s.TLSCert, s.TLSKey)
	}

	// No TLS requirements specified, listen on specified server port via http
	return httpServer.ListenAndServe()
}

// validate validates the server configuration and checks for conflicting parameters.
func (s *Server) validate() error {
	if !s.Validation {
		return nil
	}

	if s.AutoTLS && (s.TLSCert != "" || s.TLSKey != "") {
		return errors.New("AutoTLS cannot be set along with TLS cert or TLS key")
	}

	if s.AutoTLS && len(s.Domains) == 0 {
		return errors.New("AutoTLS requires a domain to also be configured")
	}

	if s.TLSCert != "" && s.TLSKey == "" {
		return errors.New("TLS certificate is missing TLS key")
	}

	if s.TLSCert == "" && s.TLSKey != "" {
		return errors.New("TLS key is missing TLS certificate")
	}

	validLogFormats := []string{"json", "text", ""}
	if !slices.Contains(validLogFormats, s.LogFormat) {
		return fmt.Errorf("invalid log format. Valid log formats are: %v", validLogFormats)
	}

	if s.LogLevel != "" {
		_, err := log.ParseLevel(s.LogLevel)
		if err != nil {
			return err
		}
	}

	return nil
}

// getLogFormatter converts a log format string to usable log formatter
func getLogFormatter(logformat string) log.Formatter {
	switch logformat {
	case "json":
		return log.JSONFormatter
	}
	return log.TextFormatter
}

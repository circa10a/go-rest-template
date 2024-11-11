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

	"github.com/caddyserver/certmagic"
	"github.com/circa10a/go-rest-template/internal/server/handlers"
	"github.com/circa10a/go-rest-template/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed api.html
var apiDocs []byte

// Server is our web server that runs the network mirror.
type Server struct {
	mux         http.Handler
	logger      *slog.Logger
	logformat   string
	loglevel    string
	tlsCert     string
	tlsKey      string
	domains     []string
	middlewares []func(http.Handler) http.Handler
	port        int
	autoTLS     bool
	metrics     bool
	validation  bool
}

// New returns a new network mirror server.
func New(options ...func(*Server)) (*Server, error) {
	server := &Server{
		loglevel: "info",
	}

	mux := http.NewServeMux()
	server.mux = mux

	// Configure server
	for _, o := range options {
		o(server)
	}

	// Ensure configuration options are valid/compatible
	err := server.validate()
	if err != nil {
		return nil, err
	}

	logLevel, err := log.ParseLevel(server.loglevel)
	if err != nil {
		return nil, err
	}

	logHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Formatter:       getLogFormatter(server.logformat),
		Level:           logLevel,
	})
	server.logger = slog.New(logHandler)

	// Features
	if server.metrics {
		mux.Handle("/metrics", promhttp.Handler())
		server.middlewares = append(server.middlewares, middleware.Prometheus)
	}

	// Default middlewares
	server.mux = middleware.Logging(server.logger, server.mux)

	// Add middlewares via http.Handler chaining
	for _, mw := range server.middlewares {
		server.mux = mw(server.mux)
	}

	// Routes
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(apiDocs) })
	mux.HandleFunc("GET /health", handlers.HealthHandleFunc)

	return server, nil
}

// Start starts the listener of the server.
func (s *Server) Start() error {
	log := s.logger.With("component", "server")

	// Auto TLS will create listeners on port 80 and 443
	if s.autoTLS {
		log.Info("Starting server on :80 and :443")
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = "user@oss.com"
		return certmagic.HTTPS(s.domains, s.mux)
	}

	// If no auto TLS, use specified server port
	// :{port}
	addr := fmt.Sprintf(":%d", s.port)
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
	if s.tlsCert != "" && s.tlsKey != "" {
		return httpServer.ListenAndServeTLS(s.tlsCert, s.tlsKey)
	}

	// No TLS requirements specified, listen on specified server port via http
	return httpServer.ListenAndServe()
}

// WithDomain configures the domain name(s) to issue a cert for with auto TLS.
func WithDomains(domains []string) func(*Server) {
	return func(s *Server) {
		s.domains = domains
	}
}

// WithMiddlewares adds middleware to the mux.
func WithMiddlewares(mws ...func(http.Handler) http.Handler) func(*Server) {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, mws...)
	}
}

// WithTLSCert configures the certificate path for a custom certificate.
func WithTLSCert(tlsCertPath string) func(*Server) {
	return func(s *Server) {
		s.tlsCert = tlsCertPath
	}
}

// WithTLSKey configures the certificate key path for a custom key.
func WithTLSKey(tlsKeyPath string) func(*Server) {
	return func(s *Server) {
		s.tlsKey = tlsKeyPath
	}
}

// WithPort configures the server to listen on the specified port when calling Start().
func WithPort(port int) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

// WithAutoTLS configures the ability to use automatic TLS or not.
func WithAutoTLS(autoTLS bool) func(*Server) {
	return func(s *Server) {
		s.autoTLS = autoTLS
	}
}

// WithMetrics configures the enablement of Prometheus metrics or not.
func WithMetrics(metrics bool) func(*Server) {
	return func(s *Server) {
		s.metrics = metrics

	}
}

// WithLogFormat configures the log format of the server.
func WithLogFormat(logformat string) func(*Server) {
	return func(s *Server) {
		s.logformat = strings.ToLower(logformat)
	}
}

// WithLogLevel configures the log level of the server.
func WithLogLevel(loglevel string) func(*Server) {
	return func(s *Server) {
		s.loglevel = loglevel
	}
}

// WithValidation ensures configuration options are valid when creating a new server.
func WithValidation(validate bool) func(*Server) {
	return func(s *Server) {
		s.validation = validate
	}
}

// validate validates the server configuration and checks for conflicting parameters.
func (s *Server) validate() error {
	if !s.validation {
		return nil
	}

	if s.autoTLS && (s.tlsCert != "" || s.tlsKey != "") {
		return errors.New("AutoTLS cannot be set along with TLS cert or TLS key")
	}

	if s.autoTLS && len(s.domains) == 0 {
		return errors.New("AutoTLS requires a domain to also be configured")
	}

	if s.tlsCert != "" && s.tlsKey == "" {
		return errors.New("TLS certificate is missing TLS key")
	}

	if s.tlsCert == "" && s.tlsKey != "" {
		return errors.New("TLS key is missing TLS certificate")
	}

	validLogFormats := []string{"json", "text", ""}
	if !slices.Contains(validLogFormats, s.logformat) {
		return fmt.Errorf("invalid log format. Valid log formats are: %v", validLogFormats)
	}

	if s.loglevel != "" {
		_, err := log.ParseLevel(s.loglevel)
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

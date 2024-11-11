package server

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		server    *Server
		expectErr bool
	}{
		{
			// No log format set
			server:    &Server{},
			expectErr: true,
		},
		{
			// No log level
			server: &Server{
				logformat: "text",
				loglevel:  "info",
			},
			expectErr: true,
		},
		{
			// Invalid log format
			server: &Server{
				logformat: "fake",
				loglevel:  "info",
			},
			expectErr: true,
		},
		{
			// Auto TLS and custom cert set
			server: &Server{
				autoTLS: true,
				tlsCert: "cert",
			},
			expectErr: true,
		},
		{
			// Auto TLS and custom key set
			server: &Server{
				autoTLS: true,
				tlsKey:  "key",
			},
			expectErr: true,
		},
		{
			// Auto TLS and no domains
			server: &Server{
				autoTLS: true,
			},
			expectErr: true,
		},
		{
			// Cert set without key
			server: &Server{
				tlsCert: "cert",
			},
			expectErr: true,
		},
		{
			// Key set without cert
			server: &Server{
				tlsKey: "key",
			},
			expectErr: true,
		},
		{
			// Valid AutoTLS config
			server: &Server{
				logformat: "text",
				loglevel:  "info",
				autoTLS:   true,
				domains:   []string{"domain"},
			},
		},
		{
			// Valid custom cert/key config
			server: &Server{
				logformat: "text",
				loglevel:  "info",
				tlsCert:   "cert",
				tlsKey:    "key",
			},
		},
	}

	for _, test := range tests {
		err := test.server.validate()
		if err != nil && !test.expectErr {
			t.Errorf("unexpected error encountered during server validation: got %s", err.Error())
		}
	}
}

func TestGetLogFormatter(t *testing.T) {
	tests := []struct {
		input    string
		expected log.Formatter
	}{
		{
			input:    "json",
			expected: log.JSONFormatter,
		},
		{
			input:    "text",
			expected: log.TextFormatter,
		},
		{
			input:    "fake",
			expected: log.TextFormatter,
		},
	}
	for _, test := range tests {
		actual := getLogFormatter(test.input)
		if test.expected != actual {
			t.Errorf("getLogFormatter returned unexpected log formatter: got %v want %v", actual, test.expected)
		}
	}
}

func TestServerConfigOpts(t *testing.T) {
	outputStr := "got: %v, want: %v"

	t.Run("WithDomains", func(t *testing.T) {
		v := []string{"lemon"}
		s, err := New(
			WithDomains(v),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if !reflect.DeepEqual(s.domains, v) {
			t.Errorf(outputStr, s.domains, v)
		}
	})

	t.Run("WithMiddlwares", func(t *testing.T) {
		v := []func(http.Handler) http.Handler{
			func(h http.Handler) http.Handler { return nil },
		}
		s, err := New(
			WithMiddlewares(v...),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if len(s.middlewares) != len(v) {
			t.Errorf(outputStr, len(s.middlewares), len(v))
		}
	})

	t.Run("WithTLSCert", func(t *testing.T) {
		v := "cert"
		s, err := New(
			WithTLSCert(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.tlsCert != v {
			t.Errorf(outputStr, s.tlsCert, v)
		}
	})

	t.Run("WithTLSKey", func(t *testing.T) {
		v := "key"
		s, err := New(
			WithTLSKey(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.tlsKey != v {
			t.Errorf(outputStr, s.tlsKey, v)
		}
	})

	t.Run("WithTLSKey", func(t *testing.T) {
		v := "key"
		s, err := New(
			WithTLSKey(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.tlsKey != v {
			t.Errorf(outputStr, s.tlsKey, v)
		}
	})

	t.Run("WithPort", func(t *testing.T) {
		v := 3000
		s, err := New(
			WithPort(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.port != v {
			t.Errorf(outputStr, s.port, v)
		}
	})

	t.Run("WithAutoTLS", func(t *testing.T) {
		v := true
		s, err := New(
			WithAutoTLS(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.autoTLS != v {
			t.Errorf(outputStr, s.autoTLS, v)
		}
	})

	t.Run("WithMetrics", func(t *testing.T) {
		v := true
		s, err := New(
			WithMetrics(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.metrics != v {
			t.Errorf(outputStr, s.metrics, v)
		}
	})

	t.Run("WithLogFormat", func(t *testing.T) {
		v := "JSON"
		vlower := strings.ToLower(v)
		s, err := New(
			WithLogFormat(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.logformat != vlower {
			t.Errorf(outputStr, s.logformat, vlower)
		}
	})

	t.Run("WithLogLevel", func(t *testing.T) {
		v := "DEBUG"
		s, err := New(
			WithLogLevel(v),
			WithValidation(false),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.loglevel != v {
			t.Errorf(outputStr, s.loglevel, v)
		}
	})

	t.Run("WithValidation", func(t *testing.T) {
		v := true
		s, err := New(
			WithValidation(v),
		)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.validation != v {
			t.Errorf(outputStr, s.validation, v)
		}
	})
}

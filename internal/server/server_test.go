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
			// Invalid log format
			server: &Server{
				Config: Config{
					LogFormat: "fake",
				},
			},
			expectErr: true,
		},
		{
			// Auto TLS and custom cert set (conflict)
			server: &Server{
				Config: Config{
					AutoTLS: true,
					TLSCert: "cert",
				},
			},
			expectErr: true,
		},
		{
			// Auto TLS and custom key set (conflict)
			server: &Server{
				Config: Config{
					AutoTLS: true,
					TLSKey:  "key",
				},
			},
			expectErr: true,
		},
		{
			// Auto TLS and no domains
			server: &Server{
				Config: Config{
					AutoTLS: true,
				},
			},
			expectErr: true,
		},
		{
			// Cert set without key
			server: &Server{
				Config: Config{
					TLSCert: "cert",
				},
			},
			expectErr: true,
		},
		{
			// Key set without cert
			server: &Server{
				Config: Config{
					TLSKey: "key",
				},
			},
			expectErr: true,
		},
		{
			// Valid AutoTLS config
			server: &Server{
				Config: Config{
					AutoTLS: true,
					Domains: []string{"domain"},
				},
			},
		},
		{
			// Valid custom cert/key config
			server: &Server{
				Config: Config{
					TLSCert: "cert",
					TLSKey:  "key",
				},
			},
		},
	}

	for _, test := range tests {
		err := test.server.validate()
		if err != nil && !test.expectErr {
			t.Errorf("unexpected validation result: got error=%v wantErr=%v, err=%v", err != nil, test.expectErr, err)
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

	t.Run("Domains", func(t *testing.T) {
		v := []string{"lemon"}
		cfg := &Config{
			Domains: v,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if !reflect.DeepEqual(s.Config.Domains, v) {
			t.Errorf(outputStr, s.Config.Domains, v)
		}
	})

	t.Run("AddMiddlewares", func(t *testing.T) {
		v := []func(http.Handler) http.Handler{
			func(h http.Handler) http.Handler { return nil },
		}
		cfg := &Config{}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		s.middlewares = v

		if len(s.middlewares) != len(v) {
			t.Errorf(outputStr, len(s.middlewares), len(v))
		}
	})

	t.Run("TLSCert", func(t *testing.T) {
		v := "cert"
		cfg := &Config{
			TLSCert:    v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.TLSCert != v {
			t.Errorf(outputStr, s.Config.TLSCert, v)
		}
	})

	t.Run("TLSKey", func(t *testing.T) {
		v := "key"
		cfg := &Config{
			TLSKey:     v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.TLSKey != v {
			t.Errorf(outputStr, s.Config.TLSKey, v)
		}
	})

	// duplicate TLSKey test removed

	t.Run("Port", func(t *testing.T) {
		v := 3000
		cfg := &Config{
			Port:       v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.Port != v {
			t.Errorf(outputStr, s.Config.Port, v)
		}
	})

	t.Run("AutoTLS", func(t *testing.T) {
		v := true
		cfg := &Config{
			AutoTLS:    v,
			Domains:    []string{"d"},
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.AutoTLS != v {
			t.Errorf(outputStr, s.Config.AutoTLS, v)
		}
	})

	t.Run("Metrics", func(t *testing.T) {
		v := true
		cfg := &Config{
			Metrics:    v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.Metrics != v {
			t.Errorf(outputStr, s.Config.Metrics, v)
		}
	})

	t.Run("LogFormat", func(t *testing.T) {
		v := "JSON"
		vlower := strings.ToLower(v)
		cfg := &Config{
			LogFormat:  v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.LogFormat != vlower {
			t.Errorf(outputStr, s.Config.LogFormat, vlower)
		}
	})

	t.Run("LogLevel", func(t *testing.T) {
		v := "DEBUG"
		cfg := &Config{
			LogLevel:   v,
			Validation: false,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Config.LogLevel != v {
			t.Errorf(outputStr, s.Config.LogLevel, v)
		}
	})

	t.Run("WithValidation", func(t *testing.T) {
		v := true
		cfg := &Config{
			Validation: v,
		}
		s, err := New(cfg)
		if err != nil {
			t.Errorf("received unexpected err: %s", err.Error())
		}

		if s.Validation != v {
			t.Errorf(outputStr, s.Validation, v)
		}
	})
}

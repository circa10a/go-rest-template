package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandleFunc(t *testing.T) {
	expected := `{"status":"ok"}`

	req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandleFunc)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %d want %d", rec.Code, http.StatusOK)
	}

	// JSON encoder adds new line endings
	actual := strings.TrimSpace(rec.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %s want %s", actual, expected)
	}
}

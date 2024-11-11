package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/circa10a/go-rest-template/api"
)

/*
HealthHandleFunc is the http handler that handles requests/responses to indicate the application is up and listening.
Example response:

	{
		"status": "ok"
	  }
*/
func HealthHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	resp := api.Health{Status: "ok"}
	_ = json.NewEncoder(w).Encode(resp)
}

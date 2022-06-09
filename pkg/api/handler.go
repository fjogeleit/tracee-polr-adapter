package api

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/fjogeleit/tracee-polr-adapter/pkg/kubernetes"
	"github.com/fjogeleit/tracee-polr-adapter/pkg/tracee"
)

// HealthzHandler for the Halthz REST API
func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, "{}")
	}
}

// ReadyHandler for the Halthz REST API
func ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{}")
	}
}

// ReadyHandler for the Halthz REST API
func WebhookHandler(polrClient *kubernetes.Client, filter *tracee.Filter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{ "message": "invalid request method - only POST is allowed" }`)
			return
		}

		var event tracee.Event

		err := json.NewDecoder(req.Body).Decode(&event)
		if err != nil {
			log.Printf("[ERROR] unable to convert event: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "message": "%s" }`, html.EscapeString(err.Error()))
			return
		}

		if s, ok := event.SigMetadata.Properties[tracee.SeverityKey]; ok {
			if severity, ok := s.(float64); ok {
				event.SigMetadata.Severity = int(severity)
			}

			delete(event.SigMetadata.Properties, tracee.SeverityKey)
		}
		if filter.Check(event) {
			err = polrClient.ProcessEvent(req.Context(), event)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{ "message": "%s" }`, html.EscapeString(err.Error()))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{}")
	}
}

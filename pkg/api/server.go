package api

import (
	"context"
	"fmt"
	"net/http"

	pprof "net/http/pprof"

	"github.com/fjogeleit/tracee-polr-adapter/pkg/kubernetes"
	"github.com/fjogeleit/tracee-polr-adapter/pkg/tracee"
)

// Server for the Lifecycle and optional HTTP REST API
type Server interface {
	// Start the HTTP Server
	Start() error
	// Shutdown the HTTP Sever
	Shutdown(ctx context.Context) error
	// RegisterLifecycleHandler adds healthy and readiness APIs
	RegisterLifecycleHandler()
	// RegisterWebhookHandler adds webhook api for tracee events
	RegisterWebhookHandler()
	// RegisterProfilingHandler adds the optional pprof profiling APIs
	RegisterProfilingHandler()
}

type httpServer struct {
	http   http.Server
	mux    *http.ServeMux
	client *kubernetes.Client
	filter *tracee.Filter
}

func (s *httpServer) RegisterLifecycleHandler() {
	s.mux.HandleFunc("/healthz", HealthzHandler())
	s.mux.HandleFunc("/ready", ReadyHandler())
}

func (s *httpServer) RegisterWebhookHandler() {
	s.mux.HandleFunc("/webhook", WebhookHandler(s.client, s.filter))
}

func (s *httpServer) RegisterProfilingHandler() {
	s.mux.HandleFunc("/debug/pprof/", pprof.Index)
	s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func (s *httpServer) Start() error {
	return s.http.ListenAndServe()
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// NewServer constructor for a new API Server
func NewServer(port int, client *kubernetes.Client, filter *tracee.Filter) Server {
	mux := http.NewServeMux()

	s := &httpServer{
		mux: mux,
		http: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		client: client,
		filter: filter,
	}

	s.RegisterLifecycleHandler()

	return s
}

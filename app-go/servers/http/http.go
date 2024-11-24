package http

import (
	"context"
	"log"
	"net/http"
	nethttp "net/http"
	"time"
)

type httpserver struct {
	server *nethttp.Server
}

type IRouter interface {
	GetRoutes() *http.ServeMux
}

// New initiate http server instance
func New(httpaddress string, handlers ...IRouter) *httpserver {
	// construct fresh mux
	mux := nethttp.NewServeMux()

	// register all handler into the server mux
	for _, handler := range handlers {
		mux.Handle("/", handler.GetRoutes())
	}

	return &httpserver{
		server: &nethttp.Server{
			Addr:    httpaddress,
			Handler: mux, // register mux into the http server
		},
	}
}

// Start starts http server
func (h *httpserver) Start() {
	go func() {
		if err := h.server.ListenAndServe(); err != nethttp.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()
}

func (h *httpserver) Stop() {
	// context for stop timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop server
	err := h.server.Shutdown(ctx)
	if err != nil {
		log.Printf("found error on stopping http server. err=%v", err)
	}
}

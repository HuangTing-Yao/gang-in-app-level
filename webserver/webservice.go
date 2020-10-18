package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var lock sync.RWMutex

type WebService struct {
	httpServer *http.Server
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, webRoute := range webRoutes {
		handler := loggingHandler(webRoute.HandlerFunc, webRoute.Name)
		router.Methods(webRoute.Method).Path(webRoute.Pattern).Name(webRoute.Name).Handler(handler)
	}
	return router
}

func loggingHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.RequestURI, name, time.Since(start))
	})
}

func (m *WebService) StartWebApp() {
	router := newRouter()
	m.httpServer = &http.Server{Addr: ":8863", Handler: router}
	log.Println("web-app started with port 8863")
	// go func() {
	httpError := m.httpServer.ListenAndServe()
	if httpError != nil && httpError != http.ErrServerClosed {
		log.Println("HTTP server error.")
	}
	// }()
}

func NewWebApp() *WebService {
	m := &WebService{}
	return m
}

func (m *WebService) StopWebApp() error {
	if m.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.httpServer.Shutdown(ctx)
	}
	return nil
}

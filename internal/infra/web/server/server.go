package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type WebServer struct {
	Router   chi.Router
	Handlers map[string]http.HandlerFunc
	Port     string
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		Port:     port,
		Router:   chi.NewRouter(),
		Handlers: make(map[string]http.HandlerFunc),
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	if err := http.ListenAndServe(s.Port, s.Router); err != nil {
		panic(err.Error())
	}
}

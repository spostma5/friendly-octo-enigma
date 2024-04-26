package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/spostma5/friendly-octo-enigma/services/risk"
)

type Server struct {
	addr string
	port int
}

func NewServer(a string, p int) *Server {
	return &Server{
		addr: a,
		port: p,
	}
}

func (s *Server) Run() error {
	slog.Info("Starting server", "addr", s.addr, "port", s.port)

	router := http.NewServeMux()
	router.HandleFunc("GET /risks", risk.HandleGetRisks)
	router.HandleFunc("POST /risks", risk.HandlePostRisk)
	router.HandleFunc("GET /risks/{id}", risk.HandleGetRisk)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.addr, s.port),
		Handler: logging(v1),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

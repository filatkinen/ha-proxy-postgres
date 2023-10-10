package backend

import (
	"context"
	"encoding/hex"
	"github.com/filatkinen/ha-proxy-postgres/internal/storage"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	db  *storage.Database
	srv *http.Server
	id  string
}

func New(port string, dbConnString string) (*Server, error) {
	serv := Server{
		srv: &http.Server{
			Addr:              net.JoinHostPort("0.0.0.0", port),
			ReadHeaderTimeout: time.Second * 5,
		},
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	serv.id = hex.EncodeToString(id[:])

	serv.srv.Handler = serv.getroute()

	db, err := storage.New(dbConnString)
	if err != nil {
		return nil, err
	}
	serv.db = db
	return &serv, nil
}

func (s *Server) Start() error {
	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("Error starting HTTP Server %s", err)
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Printf("HTTP shutdown error: %v", err)
		return err
	}
	log.Printf("HTTP graceful shutdown complete.")
	return nil
}

func (s *Server) Close() error {
	return s.db.Close()
}

func (s *Server) GetID() string {
	return s.id
}

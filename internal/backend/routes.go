package backend

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) getroute() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/getname", s.GetRandomUserName).Methods("GET")
	r.HandleFunc("/errors", s.GetErrorsDB).Methods("GET")
	r.HandleFunc("/errorscount", s.GetErrorsDBCount).Methods("GET")
	r.HandleFunc("/errorsping", s.GetErrorsPingDB).Methods("GET")
	return r
}

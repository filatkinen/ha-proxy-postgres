package backend

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) GetRandomUserName(w http.ResponseWriter, _ *http.Request) {
	name, err := s.db.SimpleQueryReturnRandomUserName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, name, "\n")
}

func (s *Server) GetErrorsDB(w http.ResponseWriter, r *http.Request) {
	errSlice := s.db.GetErrors()
	sb := strings.Builder{}
	for _, e := range errSlice {
		sb.WriteString(e)
		sb.WriteString("\n")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s.GetID()+"\n")
	fmt.Fprint(w, sb.String())
}

func (s *Server) GetErrorsDBCount(w http.ResponseWriter, r *http.Request) {
	count := s.db.GetErrorsCount()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s.GetID()+"\n")
	fmt.Fprint(w, count, "\n")
}

func (s *Server) GetErrorsPingDB(w http.ResponseWriter, r *http.Request) {
	count := s.db.GetErrorsPingDB()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s.GetID()+"\n")
	fmt.Fprint(w, count, "\n")
}

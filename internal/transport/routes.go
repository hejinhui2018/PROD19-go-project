package transport

import (
	"github.com/wogo-prod19/railguard/internal/service"
	"net/http"
)

func Routes(s *service.Service) *http.ServeMux {
	m := http.NewServeMux()
	m.Handle("/blockades/", Handler(s))
	m.HandleFunc("/recovery", func(w http.ResponseWriter, r *http.Request) {
		rep, err := s.RecoveryReport()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		WriteJSON(w, 200, rep)
	})
	m.HandleFunc("/outbox/drain", func(w http.ResponseWriter, r *http.Request) {
		rep, err := s.DrainNotices()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		WriteJSON(w, 200, rep)
	})
	return m
}

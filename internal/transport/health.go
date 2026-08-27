package transport

import (
	"github.com/wogo-prod19/railguard/internal/service"
	"net/http"
)

func HealthHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { WriteJSON(w, 200, s.Health()) }
}
func RecoveryHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep, e := s.RecoveryReport()
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		WriteJSON(w, 200, rep)
	}
}
func DrainHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rep, e := s.DrainNotices()
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		WriteJSON(w, 200, rep)
	}
}
func RegisterOperational(m *http.ServeMux, s *service.Service) {
	m.HandleFunc("/health", HealthHandler(s))
	m.HandleFunc("/recovery/report", RecoveryHandler(s))
	m.HandleFunc("/outbox/results", DrainHandler(s))
}

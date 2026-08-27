package transport

import (
	"encoding/json"
	"github.com/wogo-prod19/railguard/internal/domain"
	"github.com/wogo-prod19/railguard/internal/service"
	"net/http"
	"strings"
)

func Handler(s *service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/blockades/")
		if r.Method == "POST" && id != "" {
			var x struct {
				State domain.State `json:"state"`
			}
			if json.NewDecoder(r.Body).Decode(&x) != nil {
				http.Error(w, "bad json", 400)
				return
			}
			if e := s.Transition(id, x.State); e != nil {
				http.Error(w, e.Error(), 409)
				return
			}
			w.WriteHeader(204)
			return
		}
		b, e := s.Get(id)
		if e != nil {
			http.Error(w, e.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(b)
	})
}

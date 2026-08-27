package transport

import (
	"github.com/wogo-prod19/railguard/internal/clock"
	"github.com/wogo-prod19/railguard/internal/service"
	"net/http/httptest"
	"testing"
)

func TestMissing(t *testing.T) {
	s, _ := service.New(t.TempDir(), clock.Real{})
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/blockades/no", nil)
	Handler(s).ServeHTTP(r, req)
	if r.Code != 404 {
		t.Fatal(r.Code)
	}
}

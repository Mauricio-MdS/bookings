package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/mauricio-mds/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v:= mux.(type) {
	case *chi.Mux:
		// do nothing: expected value
	default:
		msg := fmt.Sprintf("type %T is not *chi.Mux", v)
		t.Error(msg)
	}
}
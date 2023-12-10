package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mockHandler myHandler
	handler := NoSurf(&mockHandler)

	switch v:= handler.(type) {
	case http.Handler:
		// do nothing. Expected value
	default:
		msg:= fmt.Sprintf("type %T is not http.Handler", v)
		t.Error(msg)
	}
}

func TestSessionLoad(t *testing.T) {
	var mockHandler myHandler
	handler := SessionLoad(&mockHandler)

	switch v:= handler.(type) {
	case http.Handler:
		// do nothing. Expected value
	default:
		msg:= fmt.Sprintf("type %T is not http.Handler", v)
		t.Error(msg)
	}
}
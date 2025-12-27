package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuiltinFetch(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello from test server"))
	}))
	defer server.Close()

	// Call fetch
	args := []*Expr{makeStr(server.URL)}
	result := builtinFetch(args)

	if result.Type != String {
		t.Errorf("fetch should return String, got %v", result.Type)
	}

	if result.Str != "Hello from test server" {
		t.Errorf("fetch returned %q, want 'Hello from test server'", result.Str)
	}
}

func TestBuiltinFetchJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name":"Alice","age":30}`))
	}))
	defer server.Close()

	args := []*Expr{makeStr(server.URL)}
	result := builtinFetch(args)

	expected := `{"name":"Alice","age":30}`
	if result.Str != expected {
		t.Errorf("fetch returned %q, want %q", result.Str, expected)
	}
}

func TestBuiltinFetchNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Error("fetch with 404 should panic")
		}
	}()

	args := []*Expr{makeStr(server.URL)}
	builtinFetch(args)
}

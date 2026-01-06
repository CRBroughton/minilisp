package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestFetchSimpleGet(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch (hash "url" "https://api.github.com/zen"))`
	result := eval(readStr(code), env)

	if result.Type != Hash {
		t.Fatalf("fetch result type = %v, want Hash", result.Type)
	}

	typeVal, ok := hashGet(result, "type")
	if !ok || typeVal.Str != "ok" {
		t.Error("fetch should return Ok result")
	}

	valueVal, ok := hashGet(result, "value")
	if !ok || valueVal.Type != String {
		t.Error("Ok result should have string value")
	}
}

func TestFetchWithMethod(t *testing.T) {
	env := setupTestEnv()

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		code := `(fetch (hash "url" "https://httpbin.org/anything" "method" "` + method + `"))`
		result := eval(readStr(code), env)

		typeVal, _ := hashGet(result, "type")
		if typeVal.Str != "ok" {
			t.Errorf("%s: should return Ok result", method)
		}
	}
}

func TestFetchDefaultMethod(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch (hash "url" "https://httpbin.org/get"))`
	result := eval(readStr(code), env)

	typeVal, _ := hashGet(result, "type")
	if typeVal.Str != "ok" {
		t.Error("fetch without method should default to GET")
	}
}

func TestFetchWithHeaders(t *testing.T) {
	env := setupTestEnv()

	code := `
		(fetch (hash
			"url" "https://httpbin.org/headers"
			"headers" (hash
				"User-Agent" "MiniLisp/1.0"
				"Accept" "application/json")))
	`

	result := eval(readStr(code), env)
	typeVal, _ := hashGet(result, "type")

	if typeVal.Str != "ok" {
		t.Error("fetch with headers should succeed")
	}
}

func TestFetchWithBody(t *testing.T) {
	env := setupTestEnv()

	code := `
		(fetch (hash
			"url" "https://httpbin.org/post"
			"method" "POST"
			"headers" (hash "Content-Type" "application/json")
			"body" "{\"name\": \"test\"}"))
	`

	result := eval(readStr(code), env)
	typeVal, _ := hashGet(result, "type")

	if typeVal.Str != "ok" {
		t.Error("POST with body should succeed")
	}
}

func TestFetchErrorHandling(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch (hash "url" "https://this-domain-definitely-does-not-exist-12345.com"))`
	result := eval(readStr(code), env)

	typeVal, _ := hashGet(result, "type")
	if typeVal.Str != "err" {
		t.Error("failed fetch should return Err result")
	}

	errVal, ok := hashGet(result, "error")
	if !ok || errVal.Type != String {
		t.Error("Err result should have error message")
	}
}

func TestFetchNon200Status(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch (hash "url" "https://httpbin.org/status/404"))`
	result := eval(readStr(code), env)

	typeVal, _ := hashGet(result, "type")
	if typeVal.Str != "err" {
		t.Error("404 response should return Err result")
	}

	errMsg, _ := hashGet(result, "error")
	if !strings.Contains(errMsg.Str, "404") {
		t.Error("error message should mention 404 status")
	}
}

func TestFetchInvalidArgument(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch "https://example.com")`
	result := eval(readStr(code), env)

	typeVal, _ := hashGet(result, "type")
	if typeVal.Str != "err" {
		t.Error("fetch with string should return Err result")
	}
}

func TestFetchMissingUrl(t *testing.T) {
	env := setupTestEnv()

	code := `(fetch (hash "method" "GET"))`
	result := eval(readStr(code), env)

	typeVal, _ := hashGet(result, "type")
	if typeVal.Str != "err" {
		t.Error("fetch without URL should return Err result")
	}
}

func TestFetchComplete(t *testing.T) {
	env := setupTestEnv()

	code := `
		(begin
			(define result (fetch (hash
				"url" "https://api.github.com/zen"
				"headers" (hash "User-Agent" "MiniLisp/1.0"))))

			(if (ok? result)
				(unwrap result)
				"failed"))
	`

	result := eval(readStr(code), env)

	if result.Type != String {
		t.Errorf("unwrapped result type = %v, want String", result.Type)
	}
}

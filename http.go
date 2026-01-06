package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// builtinFetch performs HTTP requests and returns Result type
// Usage: (fetch (hash "url" "..." "method" "GET" "headers" {...} "body" "..."))
func builtinFetch(args []*Expr) *Expr {
	// Validate arguments
	if len(args) != 1 {
		return resultToExpr(Err[*Expr](
			errors.New("fetch: expects 1 argument (request hash)")))
	}

	arg := args[0]
	if arg.Type != Hash {
		return resultToExpr(Err[*Expr](
			errors.New("fetch: argument must be a hash with 'url' field")))
	}

	// Extract and validate URL (required)
	urlExpr, ok := hashGet(arg, "url")
	if !ok {
		return resultToExpr(Err[*Expr](
			errors.New("fetch: request hash must have 'url' field")))
	}
	if urlExpr.Type != String {
		return resultToExpr(Err[*Expr](
			errors.New("fetch: 'url' must be a string")))
	}
	url := urlExpr.Str

	// Extract method (optional, default GET)
	method := "GET"
	if methodExpr, ok := hashGet(arg, "method"); ok {
		if methodExpr.Type != String {
			return resultToExpr(Err[*Expr](
				errors.New("fetch: 'method' must be a string")))
		}
		method = strings.ToUpper(methodExpr.Str)
	}

	// Extract headers (optional)
	headers := make(map[string]string)
	if headersExpr, ok := hashGet(arg, "headers"); ok {
		if headersExpr.Type != Hash {
			return resultToExpr(Err[*Expr](
				errors.New("fetch: 'headers' must be a hash")))
		}
		for key, val := range headersExpr.HashTable {
			if val.Type == String {
				headers[key] = val.Str
			}
		}
	}

	// Extract body (optional)
	body := ""
	if bodyExpr, ok := hashGet(arg, "body"); ok {
		if bodyExpr.Type != String {
			return resultToExpr(Err[*Expr](
				errors.New("fetch: 'body' must be a string")))
		}
		body = bodyExpr.Str
	}

	// Perform the HTTP request and wrap in Result
	return resultToExpr(performHTTPRequest(url, method, headers, body))
}

// performHTTPRequest does the actual HTTP work and returns Result[*Expr]
func performHTTPRequest(url, method string, headers map[string]string, body string) Result[*Expr] {
	// Create HTTP request
	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return Err[*Expr](fmt.Errorf("fetch: failed to create request: %w", err))
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set default User-Agent if not provided
	if _, hasUA := headers["User-Agent"]; !hasUA {
		req.Header.Set("User-Agent", "MiniLisp/1.0")
	}

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Err[*Expr](fmt.Errorf("fetch: HTTP error: %w", err))
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Err[*Expr](
			fmt.Errorf("fetch: HTTP %d: %s", resp.StatusCode, resp.Status))
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Err[*Expr](fmt.Errorf("fetch: error reading response: %w", err))
	}

	// Return Ok result with body as string
	return Ok(makeStr(string(respBody)))
}

func builtinHttpServer(args []*Expr) *Expr {
	if len(args) != 2 {
		panic("http-server: expects 2 arguments (port, handler)")
	}

	port := args[0]
	handler := args[1]

	if port.Type != Number {
		panic("http-server: port must be a number")
	}

	if handler.Type != Lambda {
		panic("http-server: handler must be a lambda")
	}

	// Create HTTP handler
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		// Build request hash for Lisp
		reqHash := makeHash()
		hashSet(reqHash, "method", makeStr(r.Method))
		hashSet(reqHash, "path", makeStr(r.URL.Path))

		// Parse query parameters
		queryHash := makeHash()
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				hashSet(queryHash, key, makeStr(values[0]))
			}
		}
		hashSet(reqHash, "query", queryHash)

		// Parse headers
		headersHash := makeHash()
		for key, values := range r.Header {
			if len(values) > 0 {
				hashSet(headersHash, key, makeStr(values[0]))
			}
		}
		hashSet(reqHash, "headers", headersHash)

		// Read body
		body, _ := io.ReadAll(r.Body)
		hashSet(reqHash, "body", makeStr(string(body)))

		// Call Lisp handler
		newEnv := NewEnv(handler.Env)
		params := handler.Params
		if params != nilExpr && params.Head != nil {
			newEnv.Define(params.Head.Sym, reqHash)
		}

		response := eval(handler.Body, newEnv)

		// Extract response fields
		if response.Type != Hash {
			panic("http-server: handler must return hash")
		}

		// Get status (default 200)
		status := 200
		if statusExpr, ok := hashGet(response, "status"); ok {
			status = statusExpr.Num
		}

		// Get headers (default empty)
		if headersExpr, ok := hashGet(response, "headers"); ok {
			if headersExpr.Type == Hash {
				for key, val := range headersExpr.HashTable {
					if val.Type == String {
						w.Header().Set(key, val.Str)
					}
				}
			}
		}

		// Get body (default empty string)
		bodyStr := ""
		if bodyExpr, ok := hashGet(response, "body"); ok {
			if bodyExpr.Type == String {
				bodyStr = bodyExpr.Str
			}
		}

		// Write response
		w.WriteHeader(status)
		w.Write([]byte(bodyStr))
	}

	// Start server
	addr := ":" + strconv.Itoa(port.Num)
	fmt.Printf("Starting server on http://localhost%s\n", addr)

	err := http.ListenAndServe(addr, http.HandlerFunc(httpHandler))
	if err != nil {
		panic(fmt.Sprintf("http-server: %v", err))
	}

	return nilExpr
}

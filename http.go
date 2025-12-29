package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func builtinFetch(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("fetch: expects 1 argument (url)")
	}

	url := args[0]
	if url.Type != String {
		panic("fetch: url must be a string")
	}

	resp, err := http.Get(url.Str)
	if err != nil {
		panic(fmt.Sprintf("fetch: HTTP error: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("fetch: HTTP %d: %s", resp.StatusCode, resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("fetch: error reading response: %v", err))
	}

	return makeStr(string(body))
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

package main

import (
	"fmt"
	"io"
	"net/http"
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

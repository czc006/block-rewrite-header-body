// Package plugin_blockpath a plugin to block a path.
package plugin_block_rewrite_header_body

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

// Config holds the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type blockRewrite struct {
	name    string
	next    http.Handler
	headers map[string]string
	body    string
}

// New creates and returns a plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &blockRewrite{
		name:    name,
		next:    next,
		headers: config.Headers,
		body:    config.Body,
	}, nil
}

func (b *blockRewrite) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.EscapedPath()
	log.Printf("blockRewrite ServeHTTP %q", currentPath)

	if len(b.headers) > 0 {
		for key, value := range b.headers {

			if key == "code" {
				code, err := strconv.Atoi(value)
				if err == nil {
					rw.WriteHeader(code)
					continue
				}
			}
			log.Printf("write header: key: %q  value: %q", key, value)
			rw.Header().Set(key, value)
		}
	}

	if len(b.body) > 0 {
		rw.Header().Del("Content-Length")
		if _, err := rw.Write([]byte(b.body)); err != nil {
			log.Printf("unable to write rewrited body: %v", err)
		}

		if flusher, ok := rw.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

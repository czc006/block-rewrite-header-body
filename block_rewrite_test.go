package plugin_block_rewrite_header_body

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestOnlyHeader(t *testing.T) {
	headers := make(map[string]string)
	headers["code"] = "204"
	headers["aaa"] = "a1"
	headers["bbb"] = "b1"

	cfg := &Config{
		Headers: headers,
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "block-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "code", "204")
	assertHeader(t, recorder, "aaa", "a1")
	assertHeader(t, recorder, "bbb", "b1")
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key string, expected string) {
	t.Helper()

	if key == "code" {
		code, err := strconv.Atoi(expected)
		if err != nil {
			t.Error(err)
		}
		if recorder.Code != code {
			t.Errorf("invalid header value: %s", recorder.Header().Get(key))
		}
		return
	}

	if recorder.Header().Get(key) != expected {
		t.Errorf("invalid header value: %s", recorder.Header().Get(key))
	}
}

func assertBody(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	t.Helper()

	if !bytes.Equal([]byte(expected), recorder.Body.Bytes()) {
		t.Errorf("got body %q, want %q", recorder.Body.Bytes(), expected)
	}
}

func TestOnlyBody(t *testing.T) {
	body := "test body"

	cfg := &Config{
		Body: body,
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "block-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "code", "200")
	assertBody(t, recorder, body)
}

func TestOnlyHeaderBody(t *testing.T) {
	body := "test body"
	headers := make(map[string]string)
	headers["code"] = "204"
	headers["aaa"] = "a1"
	headers["bbb"] = "b1"

	cfg := &Config{
		Body:    body,
		Headers: headers,
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "block-rewrite")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "code", "204")
	assertHeader(t, recorder, "aaa", "a1")
	assertHeader(t, recorder, "bbb", "b1")
	assertBody(t, recorder, body)
}

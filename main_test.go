package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCRUD(t *testing.T) {
	mux := NewMux(NewStore())

	do := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	if c := do("GET", "/healthz", "").Code; c != http.StatusOK {
		t.Fatalf("healthz: got %d", c)
	}
	if c := do("POST", "/items", `{"name":"a"}`).Code; c != http.StatusCreated {
		t.Fatalf("create: got %d", c)
	}
	if c := do("PUT", "/items/1", `{"name":"b"}`).Code; c != http.StatusOK {
		t.Fatalf("update: got %d", c)
	}
	if c := do("DELETE", "/items/1", "").Code; c != http.StatusOK {
		t.Fatalf("delete: got %d", c)
	}
	if c := do("DELETE", "/items/1", "").Code; c != http.StatusNotFound {
		t.Fatalf("delete lagi: got %d", c)
	}
	if c := do("POST", "/items", `{}`).Code; c != http.StatusBadRequest {
		t.Fatalf("validasi: got %d", c)
	}
}
package koodous

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	old := koodousBaseURL
	koodousBaseURL = srv.URL
	t.Cleanup(func() {
		koodousBaseURL = old
		srv.Close()
	})
}

func TestGetInfoDoesNotTreatDetailTagAsError(t *testing.T) {
	withServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Token test-key" {
			t.Errorf("Unexpected Authorization header: %q", got)
		}
		_, _ = io.WriteString(w, `{
			"url":"https://developer.koodous.com/apks/id",
			"id":"id",
			"created_at":"2024-01-02T03:04:05Z",
			"tags":["detail"],
			"is_detected":false,
			"rating":5,
			"is_corrupted":false,
			"trusted":true
		}`)
	})

	info, err := GetInfo("test-key", "hash")
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info == nil || len(info.Tags) != 1 || info.Tags[0] != "detail" {
		t.Fatalf("Unexpected info: %+v", info)
	}
}

func TestGetInfoNotFound(t *testing.T) {
	withServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"detail":"Not found"}`)
	})

	info, err := GetInfo("test-key", "hash")
	if err != nil {
		t.Fatalf("Expected not-found to be nil error, got %v", err)
	}
	if info != nil {
		t.Fatalf("Expected nil info, got %+v", info)
	}
}

func TestGetInfoAPIError(t *testing.T) {
	withServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"detail":"bad key"}`)
	})

	_, err := GetInfo("test-key", "hash")
	if err == nil || err.Error() != "bad key" {
		t.Fatalf("Expected API detail error, got %v", err)
	}
}

func TestGetInfoMalformedJSON(t *testing.T) {
	withServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{`)
	})

	if _, err := GetInfo("test-key", "hash"); err == nil {
		t.Fatal("Expected malformed JSON error, got nil")
	}
}

func TestGetInfoRejectsOversizedResponse(t *testing.T) {
	withServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, strings.Repeat("x", maxKoodousResponseBytes+1))
	})

	if _, err := GetInfo("test-key", "hash"); err == nil {
		t.Fatal("Expected oversized response error, got nil")
	}
}

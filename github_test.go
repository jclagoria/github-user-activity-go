package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchEvents_Success(t *testing.T) {
	events := []Event{
		{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
		{Type: EventWatch, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") != "" {
			t.Errorf("unexpected If-None-Match header on first request")
		}
		w.Header().Set("ETag", `"abc123"`)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}))
	defer server.Close()

	// Override the URL for testing
	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	got, newETag, err := FetchEvents("testuser", "")
	if err != nil {
		t.Fatalf("FetchEvents() error = %v", err)
	}
	if newETag != `"abc123"` {
		t.Errorf("FetchEvents() ETag = %q, want %q", newETag, `"abc123"`)
	}
	if len(got) != 2 {
		t.Errorf("FetchEvents() returned %d events, want 2", len(got))
	}
}

func TestFetchEvents_NotModified(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") != `"abc123"` {
			t.Errorf("expected If-None-Match = %q, got %q", `"abc123"`, r.Header.Get("If-None-Match"))
		}
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	got, newETag, err := FetchEvents("testuser", `"abc123"`)
	if err != nil {
		t.Fatalf("FetchEvents() error = %v", err)
	}
	if len(got) != 0 {
		t.Errorf("FetchEvents() returned %d events on 304, want 0", len(got))
	}
	if newETag != "" {
		t.Errorf("FetchEvents() ETag = %q on 304, want empty", newETag)
	}
}

func TestFetchEvents_UserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	_, _, err := FetchEvents("nonexistent", "")
	if err == nil {
		t.Errorf("FetchEvents() error = nil, want error for 404")
	}
}

func TestFetchEvents_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "API rate limit exceeded"})
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	_, _, err := FetchEvents("testuser", "")
	if err == nil {
		t.Errorf("FetchEvents() error = nil, want error for 403")
	}
}

func TestFilterEvents(t *testing.T) {
	events := []Event{
		{Type: EventPush},
		{Type: EventWatch},
		{Type: EventPush},
		{Type: EventFork},
	}

	t.Run("empty filter returns all", func(t *testing.T) {
		got := FilterEvents(events, "")
		if len(got) != 4 {
			t.Errorf("FilterEvents() returned %d events, want 4", len(got))
		}
	})

	t.Run("filters by type", func(t *testing.T) {
		got := FilterEvents(events, EventPush)
		if len(got) != 2 {
			t.Errorf("FilterEvents() returned %d events, want 2", len(got))
		}
	})

	t.Run("no match returns empty", func(t *testing.T) {
		got := FilterEvents(events, "NonExistentEvent")
		if len(got) != 0 {
			t.Errorf("FilterEvents() returned %d events, want 0", len(got))
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		got := FilterEvents(events, "pushevent")
		if len(got) != 2 {
			t.Errorf("FilterEvents() returned %d events for lowercase, want 2", len(got))
		}
	})
}

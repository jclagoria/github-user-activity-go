package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRun_Success(t *testing.T) {
	events := []Event{
		{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
		{Type: EventWatch, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	// Capture stdout
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run("github-activity", "testuser")

	_ = w.Close()
	os.Stdout = origStdout

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if exitCode != ExitOK {
		t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
	}
	if !strings.Contains(output, "user/repo: Pushed commits") {
		t.Errorf("output missing push event: %s", output)
	}
	if !strings.Contains(output, "user/repo: Starred") {
		t.Errorf("output missing watch event: %s", output)
	}
}

func TestRun_NoArgs(t *testing.T) {
	origStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run("github-activity")

	_ = w.Close()
	os.Stderr = origStderr

	if exitCode != ExitUsage {
		t.Errorf("exit code = %d, want %d", exitCode, ExitUsage)
	}
}

func TestRun_UserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	origStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run("github-activity", "nonexistent")

	_ = w.Close()
	os.Stderr = origStderr

	if exitCode != ExitNotFound {
		t.Errorf("exit code = %d, want %d", exitCode, ExitNotFound)
	}
}

func TestRun_EmptyEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Event{})
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	origStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run("github-activity", "emptyuser")

	_ = w.Close()
	os.Stderr = origStderr

	if exitCode != ExitOK {
		t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
	}
}

func TestRun_WithFilter(t *testing.T) {
	events := []Event{
		{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
		{Type: EventWatch, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run("github-activity", "testuser", "--type", "PushEvent")

	_ = w.Close()
	os.Stdout = origStdout

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if exitCode != ExitOK {
		t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
	}
	if !strings.Contains(output, "Pushed commits") {
		t.Errorf("output missing push event: %s", output)
	}
	if strings.Contains(output, "Starred") {
		t.Errorf("output should not contain watch event with filter: %s", output)
	}
}

func TestRun_MaxEvents(t *testing.T) {
	events := make([]Event, 15)
	for i := range events {
		events[i] = Event{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}))
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run("github-activity", "testuser")

	_ = w.Close()
	os.Stdout = origStdout

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if exitCode != ExitOK {
		t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != maxEvents {
		t.Errorf("output has %d lines, want %d", len(lines), maxEvents)
	}
}

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

// withTestServer sets up an httptest server, overrides eventsURL, and calls fn.
func withTestServer(t *testing.T, handler http.HandlerFunc, fn func(serverURL string)) {
	t.Helper()
	server := httptest.NewServer(handler)
	defer server.Close()

	origURL := eventsURL
	eventsURL = server.URL + "/users/%s/events"
	defer func() { eventsURL = origURL }()

	fn(server.URL)
}

// captureStdout runs fn while capturing stdout, returns the captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = orig

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// suppressStderr runs fn while suppressing stderr output.
func suppressStderr(t *testing.T, fn func()) {
	t.Helper()
	orig := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = orig
}

func TestRun_Success(t *testing.T) {
	events := []Event{
		{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
		{Type: EventWatch, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
	}

	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}, func(_ string) {
		output := captureStdout(t, func() {
			exitCode := run("github-activity", "testuser")
			if exitCode != ExitOK {
				t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
			}
		})
		if !strings.Contains(output, "Pushed") && !strings.Contains(output, "to user/repo") {
			t.Errorf("output missing push event: %s", output)
		}
		if !strings.Contains(output, "Starred user/repo") {
			t.Errorf("output missing watch event: %s", output)
		}
	})
}

func TestRun_NoArgs(t *testing.T) {
	suppressStderr(t, func() {
		exitCode := run("github-activity")
		if exitCode != ExitUsage {
			t.Errorf("exit code = %d, want %d", exitCode, ExitUsage)
		}
	})
}

func TestRun_UserNotFound(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
	}, func(_ string) {
		suppressStderr(t, func() {
			exitCode := run("github-activity", "nonexistent")
			if exitCode != ExitNotFound {
				t.Errorf("exit code = %d, want %d", exitCode, ExitNotFound)
			}
		})
	})
}

func TestRun_EmptyEvents(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Event{})
	}, func(_ string) {
		suppressStderr(t, func() {
			exitCode := run("github-activity", "emptyuser")
			if exitCode != ExitOK {
				t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
			}
		})
	})
}

func TestRun_WithFilter(t *testing.T) {
	events := []Event{
		{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
		{Type: EventWatch, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}},
	}

	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}, func(_ string) {
		output := captureStdout(t, func() {
			exitCode := run("github-activity", "testuser", "--type", "PushEvent")
			if exitCode != ExitOK {
				t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
			}
		})
		if !strings.Contains(output, "Pushed") {
			t.Errorf("output missing push event: %s", output)
		}
		if strings.Contains(output, "Starred") {
			t.Errorf("output should not contain watch event with filter: %s", output)
		}
	})
}

func TestRun_MaxEvents(t *testing.T) {
	events := make([]Event, 15)
	for i := range events {
		events[i] = Event{Type: EventPush, Repo: Repo{Name: "user/repo"}, Actor: Actor{Login: "user"}}
	}

	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}, func(_ string) {
		output := captureStdout(t, func() {
			exitCode := run("github-activity", "testuser")
			if exitCode != ExitOK {
				t.Errorf("exit code = %d, want %d", exitCode, ExitOK)
			}
		})
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != maxEvents {
			t.Errorf("output has %d lines, want %d", len(lines), maxEvents)
		}
	})
}

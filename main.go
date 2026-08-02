package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxEvents = 10

// ParseArgs extracts the username and optional --type filter from CLI arguments.
func ParseArgs(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", fmt.Errorf("usage: github-activity <username>")
	}

	username := args[1]

	if len(args) >= 4 && args[2] == "--type" {
		return username, args[3], nil
	}

	if len(args) >= 3 && args[2] == "--type" {
		return "", "", fmt.Errorf("usage: github-activity <username> --type <EventType>")
	}

	// Ignore unknown flags, just return username
	for i := 2; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			return "", "", fmt.Errorf("unknown flag: %s", args[i])
		}
	}

	return username, "", nil
}

func main() {
	os.Exit(run(os.Args...))
}

func run(args ...string) int {
	username, typeFilter, err := ParseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, MsgUsage)
		return ExitUsage
	}

	etag := loadETag(username)
	events, newETag, err := FetchEvents(username, etag)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "404"):
			fmt.Fprintf(os.Stderr, MsgNotFound+"\n", username)
			return ExitNotFound
		case strings.Contains(msg, "403"):
			fmt.Fprintln(os.Stderr, MsgRateLimit)
			return ExitRateLimit
		case strings.Contains(msg, "failed to connect"):
			fmt.Fprintln(os.Stderr, MsgNetwork)
			return ExitNetwork
		case strings.Contains(msg, "unexpected"):
			fmt.Fprintln(os.Stderr, MsgInvalidJSON)
			return ExitInvalidJSON
		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
			return 1
		}
	}

	if newETag != "" {
		saveETag(username, newETag)
	}

	events = FilterEvents(events, typeFilter)

	if len(events) == 0 {
		fmt.Fprintf(os.Stderr, MsgNoActivity+"\n", username)
		return ExitOK
	}

	if len(events) > maxEvents {
		events = events[:maxEvents]
	}

	for _, event := range events {
		fmt.Println(FormatEvent(event))
	}

	return ExitOK
}

func etagPath(username string) string {
	dir := os.TempDir()
	return filepath.Join(dir, ".github-activity-"+username+".etag")
}

func loadETag(username string) string {
	data, err := os.ReadFile(etagPath(username))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveETag(username, etag string) {
	_ = os.WriteFile(etagPath(username), []byte(etag), 0600)
}

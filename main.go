package main

import (
	"fmt"
	"strings"
)

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
}

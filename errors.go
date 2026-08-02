package main

import "os"

// Exit codes.
const (
	ExitOK            = 0
	ExitUsage         = 1
	ExitNotFound      = 2
	ExitRateLimit     = 3
	ExitNetwork       = 4
	ExitInvalidJSON   = 5
)

// Error messages written to stderr.
var (
	MsgUsage        = "Usage: github-activity <username>"
	MsgNotFound     = "Error: user '%s' not found"
	MsgRateLimit    = "Error: API rate limit exceeded. Try again later."
	MsgNetwork      = "Error: failed to connect. Check your internet."
	MsgInvalidJSON  = "Error: unexpected API response"
	MsgNoActivity   = "No recent activity for '%s'"
)

func usageError() {
	os.Stderr.WriteString(MsgUsage + "\n")
	os.Exit(ExitUsage)
}

func notFoundError(username string) {
	os.Stderr.WriteString(MsgNotFound + "\n")
	os.Exit(ExitNotFound)
}

func rateLimitError() {
	os.Stderr.WriteString(MsgRateLimit + "\n")
	os.Exit(ExitRateLimit)
}

func networkError() {
	os.Stderr.WriteString(MsgNetwork + "\n")
	os.Exit(ExitNetwork)
}

func invalidJSONError() {
	os.Stderr.WriteString(MsgInvalidJSON + "\n")
	os.Exit(ExitInvalidJSON)
}

package main

import "fmt"

// Exit codes.
const (
	ExitOK       = 0
	ExitUsage    = 1
	ExitNotFound = 2
	ExitRateLimit = 3
	ExitNetwork  = 4
	ExitInvalidJSON = 5
)

// Typed errors returned by FetchEvents.
type NotFoundError struct{ User string }
type RateLimitError struct{}
type NetworkError struct{ Msg string }
type InvalidJSONError struct{ Msg string }

func (e *NotFoundError) Error() string  { return fmt.Sprintf("user '%s' not found", e.User) }
func (e *RateLimitError) Error() string { return "API rate limit exceeded" }
func (e *NetworkError) Error() string   { return e.Msg }
func (e *InvalidJSONError) Error() string { return e.Msg }

// Error messages written to stderr.
var (
	MsgUsage      = "Usage: github-activity <username>"
	MsgNotFound   = "Error: user '%s' not found"
	MsgRateLimit  = "Error: API rate limit exceeded. Try again later."
	MsgNetwork    = "Error: failed to connect. Check your internet."
	MsgInvalidJSON = "Error: unexpected API response"
	MsgNoActivity = "No recent activity for '%s'"
)

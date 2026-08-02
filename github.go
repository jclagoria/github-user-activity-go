package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Event represents a GitHub event from the Events API.
type Event struct {
	Type    string          `json:"type"`
	Repo    Repo            `json:"repo"`
	Actor   Actor           `json:"actor"`
	Payload json.RawMessage `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
}

type Actor struct {
	Login string `json:"login"`
}

// Event type constants.
const (
	EventPush                     = "PushEvent"
	EventIssues                   = "IssuesEvent"
	EventIssueComment             = "IssueCommentEvent"
	EventWatch                    = "WatchEvent"
	EventFork                     = "ForkEvent"
	EventCreate                   = "CreateEvent"
	EventDelete                   = "DeleteEvent"
	EventPullRequest              = "PullRequestEvent"
	EventPullRequestReview        = "PullRequestReviewEvent"
	EventPullRequestReviewComment = "PullRequestReviewCommentEvent"
	EventRelease                  = "ReleaseEvent"
	EventCommitComment            = "CommitCommentEvent"
	EventMember                   = "MemberEvent"
	EventGollum                   = "GollumEvent"
	EventPublic                   = "PublicEvent"
	EventDiscussion               = "DiscussionEvent"
)

var eventsURL = "https://api.github.com/users/%s/events"

// FetchEvents fetches recent public events for a GitHub user.
// It supports ETag-based caching: pass the previous ETag to avoid
// re-downloading unchanged data (GitHub returns 304 Not Modified).
func FetchEvents(username string, etag string) ([]Event, string, error) {
	url := fmt.Sprintf(eventsURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	newETag := resp.Header.Get("ETag")

	if resp.StatusCode == http.StatusNotModified {
		return nil, newETag, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
			return nil, newETag, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Message)
		}
		return nil, newETag, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, newETag, fmt.Errorf("unexpected API response: %w", err)
	}

	// Normalize: GitHub sometimes returns null for events
	if events == nil {
		events = []Event{}
	}

	return events, newETag, nil
}

// FilterEvents returns only events matching the given type.
// If eventType is empty, all events are returned.
func FilterEvents(events []Event, eventType string) []Event {
	if eventType == "" {
		return events
	}

	var filtered []Event
	for _, e := range events {
		if strings.EqualFold(e.Type, eventType) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

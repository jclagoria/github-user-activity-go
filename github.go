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

// FormatEvent formats a single event into a human-readable string.
// Format: "owner/repo: Verb details"
func FormatEvent(event Event) string {
	repo := event.Repo.Name

	switch event.Type {
	case EventPush:
		return fmt.Sprintf("%s: Pushed commits", repo)
	case EventIssues:
		return formatActionVerb(event, "issue")
	case EventIssueComment:
		return fmt.Sprintf("%s: Commented on issue", repo)
	case EventWatch:
		return fmt.Sprintf("%s: Starred", repo)
	case EventFork:
		return fmt.Sprintf("%s: Forked", repo)
	case EventCreate:
		return formatCreateEvent(event)
	case EventDelete:
		return formatDeleteEvent(event)
	case EventPullRequest:
		return formatPREvent(event, "PR")
	case EventPullRequestReview:
		return fmt.Sprintf("%s: Reviewed PR", repo)
	case EventPullRequestReviewComment:
		return fmt.Sprintf("%s: Commented on PR", repo)
	case EventRelease:
		return formatReleaseEvent(event)
	case EventCommitComment:
		return fmt.Sprintf("%s: Commented on commit", repo)
	case EventMember:
		return fmt.Sprintf("%s: Added member", repo)
	case EventGollum:
		return fmt.Sprintf("%s: Updated wiki", repo)
	case EventPublic:
		return fmt.Sprintf("%s: Made repo public", repo)
	case EventDiscussion:
		return fmt.Sprintf("%s: Updated discussion", repo)
	default:
		return fmt.Sprintf("%s: Did something", repo)
	}
}

func formatActionVerb(event Event, noun string) string {
	var payload struct {
		Action string `json:"action"`
	}
	json.Unmarshal(event.Payload, &payload)

	action := strings.Title(payload.Action)
	if action == "" {
		action = "Updated"
	}
	return fmt.Sprintf("%s: %s %s", event.Repo.Name, action, noun)
}

func formatCreateEvent(event Event) string {
	var payload struct {
		RefType string `json:"ref_type"`
		Ref     string `json:"ref"`
	}
	json.Unmarshal(event.Payload, &payload)

	switch payload.RefType {
	case "repository":
		return fmt.Sprintf("%s: Created repo", event.Repo.Name)
	case "branch":
		return fmt.Sprintf("%s: Created branch %s", event.Repo.Name, payload.Ref)
	case "tag":
		return fmt.Sprintf("%s: Created tag %s", event.Repo.Name, payload.Ref)
	default:
		return fmt.Sprintf("%s: Created %s", event.Repo.Name, payload.RefType)
	}
}

func formatDeleteEvent(event Event) string {
	var payload struct {
		RefType string `json:"ref_type"`
		Ref     string `json:"ref"`
	}
	json.Unmarshal(event.Payload, &payload)

	return fmt.Sprintf("%s: Deleted %s %s", event.Repo.Name, payload.RefType, payload.Ref)
}

func formatPREvent(event Event, noun string) string {
	var payload struct {
		Action string `json:"action"`
	}
	json.Unmarshal(event.Payload, &payload)

	action := strings.Title(payload.Action)
	if action == "" {
		action = "Updated"
	}
	return fmt.Sprintf("%s: %s %s", event.Repo.Name, action, noun)
}

func formatReleaseEvent(event Event) string {
	var payload struct {
		Action  string `json:"action"`
		Release struct {
			TagName string `json:"tag_name"`
		} `json:"release"`
	}
	json.Unmarshal(event.Payload, &payload)

	action := strings.Title(payload.Action)
	if action == "" {
		action = "Updated"
	}
	return fmt.Sprintf("%s: %s release %s", event.Repo.Name, action, payload.Release.TagName)
}

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
func FetchEvents(username, etag string) ([]Event, string, error) {
	url := fmt.Sprintf(eventsURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", &NetworkError{Msg: err.Error()}
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", &NetworkError{Msg: err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	newETag := resp.Header.Get("ETag")

	if resp.StatusCode == http.StatusNotModified {
		return nil, newETag, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, newETag, &NotFoundError{User: username}
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, newETag, &RateLimitError{}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
			return nil, newETag, &InvalidJSONError{Msg: errResp.Message}
		}
		return nil, newETag, &InvalidJSONError{Msg: fmt.Sprintf("unexpected status: %d", resp.StatusCode)}
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, newETag, &InvalidJSONError{Msg: err.Error()}
	}

	if events == nil {
		events = []Event{}
	}

	return events, newETag, nil
}

// FilterEvents returns only events matching the given type.
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
// Format: "Verb ... in/to repo" per spec.
func FormatEvent(event Event) string {
	repo := event.Repo.Name

	switch event.Type {
	case EventPush:
		return formatPushEvent(event)
	case EventIssues:
		return formatActionEvent(event, "issue")
	case EventIssueComment:
		return fmt.Sprintf("Commented on issue in %s", repo)
	case EventWatch:
		return fmt.Sprintf("Starred %s", repo)
	case EventFork:
		return fmt.Sprintf("Forked %s", repo)
	case EventCreate:
		return formatCreateEvent(event)
	case EventDelete:
		return formatDeleteEvent(event)
	case EventPullRequest:
		return formatActionEvent(event, "PR")
	case EventPullRequestReview:
		return fmt.Sprintf("Reviewed PR in %s", repo)
	case EventPullRequestReviewComment:
		return fmt.Sprintf("Commented on PR in %s", repo)
	case EventRelease:
		return formatReleaseEvent(event)
	case EventCommitComment:
		return fmt.Sprintf("Commented on commit in %s", repo)
	case EventMember:
		return fmt.Sprintf("Added member in %s", repo)
	case EventGollum:
		return fmt.Sprintf("Updated wiki in %s", repo)
	case EventPublic:
		return fmt.Sprintf("Made %s public", repo)
	case EventDiscussion:
		return fmt.Sprintf("Updated discussion in %s", repo)
	default:
		return fmt.Sprintf("Did something in %s", repo)
	}
}

func formatPushEvent(event Event) string {
	var payload struct {
		Commits []struct{} `json:"commits"`
	}
	_ = json.Unmarshal(event.Payload, &payload)
	count := len(payload.Commits)
	if count == 0 {
		return fmt.Sprintf("Pushed commits to %s", event.Repo.Name)
	}
	return fmt.Sprintf("Pushed %d commits to %s", count, event.Repo.Name)
}

func formatActionEvent(event Event, noun string) string {
	var payload struct {
		Action string `json:"action"`
	}
	_ = json.Unmarshal(event.Payload, &payload)

	action := capitalize(payload.Action)
	if action == "" {
		action = "Updated"
	}

	if action == "Opened" {
		return fmt.Sprintf("Opened a new %s in %s", noun, event.Repo.Name)
	}
	return fmt.Sprintf("%s %s in %s", action, noun, event.Repo.Name)
}

func formatCreateEvent(event Event) string {
	var payload struct {
		RefType string `json:"ref_type"`
		Ref     string `json:"ref"`
	}
	_ = json.Unmarshal(event.Payload, &payload)

	switch payload.RefType {
	case "repository":
		return fmt.Sprintf("Created repo in %s", event.Repo.Name)
	case "branch":
		return fmt.Sprintf("Created branch %s in %s", payload.Ref, event.Repo.Name)
	case "tag":
		return fmt.Sprintf("Created tag %s in %s", payload.Ref, event.Repo.Name)
	default:
		return fmt.Sprintf("Created %s in %s", payload.RefType, event.Repo.Name)
	}
}

func formatDeleteEvent(event Event) string {
	var payload struct {
		RefType string `json:"ref_type"`
		Ref     string `json:"ref"`
	}
	_ = json.Unmarshal(event.Payload, &payload)

	return fmt.Sprintf("Deleted %s %s in %s", payload.RefType, payload.Ref, event.Repo.Name)
}

func formatReleaseEvent(event Event) string {
	var payload struct {
		Action  string `json:"action"`
		Release struct {
			TagName string `json:"tag_name"`
		} `json:"release"`
	}
	_ = json.Unmarshal(event.Payload, &payload)

	action := capitalize(payload.Action)
	if action == "" {
		action = "Updated"
	}
	return fmt.Sprintf("%s release %s in %s", action, payload.Release.TagName, event.Repo.Name)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

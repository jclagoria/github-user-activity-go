package main

import "encoding/json"

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
	EventPush                    = "PushEvent"
	EventIssues                  = "IssuesEvent"
	EventIssueComment            = "IssueCommentEvent"
	EventWatch                   = "WatchEvent"
	EventFork                    = "ForkEvent"
	EventCreate                  = "CreateEvent"
	EventDelete                  = "DeleteEvent"
	EventPullRequest             = "PullRequestEvent"
	EventPullRequestReview       = "PullRequestReviewEvent"
	EventPullRequestReviewComment = "PullRequestReviewCommentEvent"
	EventRelease                 = "ReleaseEvent"
	EventCommitComment           = "CommitCommentEvent"
	EventMember                  = "MemberEvent"
	EventGollum                  = "GollumEvent"
	EventPublic                  = "PublicEvent"
	EventDiscussion              = "DiscussionEvent"
)

package main

import (
	"encoding/json"
	"testing"
)

func TestFormatEvent(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		want  string
	}{
		{
			name:  "push event",
			event: Event{Type: EventPush, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Pushed commits",
		},
		{
			name:  "issues event opened",
			event: Event{Type: EventIssues, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "opened"})},
			want:  "user/repo: Opened issue",
		},
		{
			name:  "issues event closed",
			event: Event{Type: EventIssues, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "closed"})},
			want:  "user/repo: Closed issue",
		},
		{
			name:  "issue comment event",
			event: Event{Type: EventIssueComment, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Commented on issue",
		},
		{
			name:  "watch event",
			event: Event{Type: EventWatch, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Starred",
		},
		{
			name:  "fork event",
			event: Event{Type: EventFork, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Forked",
		},
		{
			name:  "create event repository",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "repository"})},
			want:  "user/repo: Created repo",
		},
		{
			name:  "create event branch",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "branch", "ref": "main"})},
			want:  "user/repo: Created branch main",
		},
		{
			name:  "create event tag",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "tag", "ref": "v1.0"})},
			want:  "user/repo: Created tag v1.0",
		},
		{
			name:  "delete event",
			event: Event{Type: EventDelete, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "branch", "ref": "feature-x"})},
			want:  "user/repo: Deleted branch feature-x",
		},
		{
			name:  "pull request event opened",
			event: Event{Type: EventPullRequest, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "opened"})},
			want:  "user/repo: Opened PR",
		},
		{
			name:  "pull request review event",
			event: Event{Type: EventPullRequestReview, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Reviewed PR",
		},
		{
			name:  "pull request review comment event",
			event: Event{Type: EventPullRequestReviewComment, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Commented on PR",
		},
		{
			name: "release event published",
			event: Event{
				Type: EventRelease,
				Repo: Repo{Name: "user/repo"},
				Payload: mustMarshal(map[string]interface{}{
					"action":  "published",
					"release": map[string]string{"tag_name": "v1.0"},
				}),
			},
			want: "user/repo: Published release v1.0",
		},
		{
			name:  "commit comment event",
			event: Event{Type: EventCommitComment, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Commented on commit",
		},
		{
			name:  "member event",
			event: Event{Type: EventMember, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Added member",
		},
		{
			name:  "gollum event",
			event: Event{Type: EventGollum, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Updated wiki",
		},
		{
			name:  "public event",
			event: Event{Type: EventPublic, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Made repo public",
		},
		{
			name:  "discussion event",
			event: Event{Type: EventDiscussion, Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Updated discussion",
		},
		{
			name:  "unknown event type",
			event: Event{Type: "UnknownEvent", Repo: Repo{Name: "user/repo"}},
			want:  "user/repo: Did something",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatEvent(tt.event)
			if got != tt.want {
				t.Errorf("FormatEvent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

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
			event: Event{Type: EventPush, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]interface{}{"commits": []int{1, 2, 3}})},
			want:  "Pushed 3 commits to user/repo",
		},
		{
			name:  "push event no commits",
			event: Event{Type: EventPush, Repo: Repo{Name: "user/repo"}},
			want:  "Pushed commits to user/repo",
		},
		{
			name:  "issues event opened",
			event: Event{Type: EventIssues, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "opened"})},
			want:  "Opened a new issue in user/repo",
		},
		{
			name:  "issues event closed",
			event: Event{Type: EventIssues, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "closed"})},
			want:  "Closed issue in user/repo",
		},
		{
			name:  "issue comment event",
			event: Event{Type: EventIssueComment, Repo: Repo{Name: "user/repo"}},
			want:  "Commented on issue in user/repo",
		},
		{
			name:  "watch event",
			event: Event{Type: EventWatch, Repo: Repo{Name: "user/repo"}},
			want:  "Starred user/repo",
		},
		{
			name:  "fork event",
			event: Event{Type: EventFork, Repo: Repo{Name: "user/repo"}},
			want:  "Forked user/repo",
		},
		{
			name:  "create event repository",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "repository"})},
			want:  "Created repo in user/repo",
		},
		{
			name:  "create event branch",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "branch", "ref": "main"})},
			want:  "Created branch main in user/repo",
		},
		{
			name:  "create event tag",
			event: Event{Type: EventCreate, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "tag", "ref": "v1.0"})},
			want:  "Created tag v1.0 in user/repo",
		},
		{
			name:  "delete event",
			event: Event{Type: EventDelete, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"ref_type": "branch", "ref": "feature-x"})},
			want:  "Deleted branch feature-x in user/repo",
		},
		{
			name:  "pull request event opened",
			event: Event{Type: EventPullRequest, Repo: Repo{Name: "user/repo"}, Payload: mustMarshal(map[string]string{"action": "opened"})},
			want:  "Opened a new PR in user/repo",
		},
		{
			name:  "pull request review event",
			event: Event{Type: EventPullRequestReview, Repo: Repo{Name: "user/repo"}},
			want:  "Reviewed PR in user/repo",
		},
		{
			name:  "pull request review comment event",
			event: Event{Type: EventPullRequestReviewComment, Repo: Repo{Name: "user/repo"}},
			want:  "Commented on PR in user/repo",
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
			want: "Published release v1.0 in user/repo",
		},
		{
			name:  "commit comment event",
			event: Event{Type: EventCommitComment, Repo: Repo{Name: "user/repo"}},
			want:  "Commented on commit in user/repo",
		},
		{
			name:  "member event",
			event: Event{Type: EventMember, Repo: Repo{Name: "user/repo"}},
			want:  "Added member in user/repo",
		},
		{
			name:  "gollum event",
			event: Event{Type: EventGollum, Repo: Repo{Name: "user/repo"}},
			want:  "Updated wiki in user/repo",
		},
		{
			name:  "public event",
			event: Event{Type: EventPublic, Repo: Repo{Name: "user/repo"}},
			want:  "Made user/repo public",
		},
		{
			name:  "discussion event",
			event: Event{Type: EventDiscussion, Repo: Repo{Name: "user/repo"}},
			want:  "Updated discussion in user/repo",
		},
		{
			name:  "unknown event type",
			event: Event{Type: "UnknownEvent", Repo: Repo{Name: "user/repo"}},
			want:  "Did something in user/repo",
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

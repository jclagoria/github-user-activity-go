# GitHub User Activity CLI

Fetch a GitHub user's recent activity and display it in the terminal.

Built as a learning project following the [GitHub User Activity](https://roadmap.sh/projects/github-user-activity) roadmap challenge.

## Install

```bash
go install github.com/jclagoria/github-user-activity-go@latest
```

Or build from source:

```bash
git clone https://github.com/jclagoria/github-user-activity-go.git
cd github-user-activity-go
go build -o github-activity .
```

## Usage

```bash
github-activity <username>
```

Example:

```bash
github-activity kamranahmedse
```

Output:

```
Pushed 3 commits to kamranahmedse/developer-roadmap
Opened a new issue in kamranahmedse/developer-roadmap
Starred kamranahmedse/developer-roadmap
```

### Filter by event type

```bash
github-activity <username> --type <EventType>
```

Example:

```bash
github-activity kamranahmedse --type PushEvent
```

Supported event types: `PushEvent`, `IssuesEvent`, `IssueCommentEvent`, `WatchEvent`, `ForkEvent`, `CreateEvent`, `DeleteEvent`, `PullRequestEvent`, `PullRequestReviewEvent`, `PullRequestReviewCommentEvent`, `ReleaseEvent`, `CommitCommentEvent`, `MemberEvent`, `GollumEvent`, `PublicEvent`, `DiscussionEvent`.

### Caching

The CLI uses ETag-based caching. On subsequent runs for the same user, unchanged responses are skipped (HTTP 304), reducing API calls and staying within the 60 req/hr unauthenticated rate limit.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (or no activity) |
| 1 | Usage error |
| 2 | User not found |
| 3 | Rate limit exceeded |
| 4 | Network error |
| 5 | Invalid API response |

## Running tests

```bash
go test ./...
```

## License

MIT

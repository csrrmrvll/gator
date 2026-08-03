# gator

`gator` is a command-line RSS feed aggregator written in Go. It lets multiple
users register, subscribe to RSS feeds, continuously fetch new posts into a
PostgreSQL database, and browse the latest entries from feeds they follow.

## Features

- Multi-user support with a "current user" stored in a local config file.
- Add and list RSS feeds shared across all users.
- Follow / unfollow feeds per user.
- A long-running aggregator that periodically fetches feeds and saves posts.
- Browse recent posts from the feeds the current user follows.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+ (see [go.mod](go.mod)).
- [PostgreSQL](https://www.postgresql.org/) 14+ running and reachable.
- (Optional, for development) [goose](https://github.com/pressly/goose) for
  running migrations and [sqlc](https://sqlc.dev/) for regenerating the DB layer.

## Installation

Install the CLI with `go install`:

```bash
go install github.com/csrrmrvll/gator@latest
```

Or clone and build from source:

```bash
git clone https://github.com/csrrmrvll/gator.git
cd gator
go build -o gator .
```

## Configuration

`gator` reads a JSON config file named `.gatorconfig.json` from your home
directory (`$HOME/.gatorconfig.json`). Create it before running any commands:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

| Field               | Description                                             |
| ------------------- | ------------------------------------------------------- |
| `db_url`            | PostgreSQL connection string.                           |
| `current_user_name` | Set automatically by `register` / `login`. Leave blank. |

## Database setup

Create the database, then apply the migrations in [sql/schema](sql/schema) using
`goose`:

```bash
createdb gator
cd sql/schema
goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

## Usage

```bash
gator <command> [args...]
```

### Commands

| Command                     | Auth | Description                                                     |
| --------------------------- | :--: | -------------------------------------------------------------- |
| `register <name>`           |  no  | Create a new user and set it as the current user.             |
| `login <name>`              |  no  | Switch the current user to an existing user.                  |
| `users`                     |  no  | List all users; the current user is marked `(current)`.       |
| `reset`                     |  no  | Delete all users (and cascading data). Destructive.          |
| `addfeed <name> <url>`      | yes  | Add a feed and automatically follow it.                      |
| `feeds`                     |  no  | List all feeds with their owner.                             |
| `follow <feed_url>`         | yes  | Follow an existing feed by URL.                              |
| `unfollow <feed_url>`       | yes  | Unfollow a feed by URL.                                      |
| `following`                 | yes  | List feeds the current user follows.                        |
| `agg <time_between_reqs>`   | yes  | Run the aggregator loop, fetching feeds on the given interval (e.g. `1m`, `30s`). |
| `browse [limit]`            | yes  | Show recent posts from followed feeds (default limit: 2).   |

Commands marked **Auth = yes** require a current user to be set via `register`
or `login`.

### Example workflow

```bash
# Create a user (becomes the current user)
gator register alice

# Add and automatically follow a feed
gator addfeed "TechCrunch" "https://techcrunch.com/feed/"

# Start the aggregator (Ctrl+C to stop) — fetch every minute
gator agg 1m

# In another shell: browse the latest 5 posts
gator browse 5
```

## Project structure

```text
main.go               Entry point, DB connection, command registration
commands.go           Command registry and dispatch
handler_*.go          One file per command group
rss_feed.go           RSS fetching and XML parsing
internal/config/      Read/write the .gatorconfig.json file
internal/database/    sqlc-generated database access layer
sql/schema/           goose migrations
sql/queries/          sqlc query definitions
sqlc.yaml             sqlc configuration
```

## Development

Regenerate the database layer after editing files in `sql/queries` or
`sql/schema`:

```bash
sqlc generate
```

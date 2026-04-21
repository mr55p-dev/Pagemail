# Pagemail — Project Context

## What It Is

Pagemail (pagemail.io) is a link-saving / read-later tool. Users save URLs, the app fetches readable content, and they can browse/manage their saved pages.

## Current State

- **Production version:** v5.4.0 (tagged)
- **Latest master:** one commit ahead — `fix/remove the articles and readings tables (#160)` — this is in production per Ellis
- **Version history:** v5.x series, currently on v5.4.x

### Recent significant changes
- `#159` — Replatformed from (unknown) to **Hetzner** for hosting
- `#157` — Added articles and readings tables (experimental feature)
- `#160` — **Reverted** articles/readings tables (they were removed, not yet in a release tag)
- `#155` — Readability synth endpoint added
- `#147` — Articles page added (now reverted)

## Tech Stack

- **Language:** Go 1.22
- **Web framework:** stdlib `net/http` with custom router
- **Templating:** [templ](https://templ.guide) (`.templ` files compiled to Go)
- **Database:** SQLite via [sqlc](https://sqlc.dev) (typed query generation)
- **Auth:** Google OAuth + email/password login
- **Deployment:** Docker Compose → AWS ECR images → Hetzner server behind Traefik
- **Sidecar:** Python `readability` service (port 5000) for article content extraction
- **AWS:** Used for config/credentials (likely SES for email, ECR for images)

## Project Structure

```
cmd/
  pagemail/        — main entrypoint
  mailmock/        — mock mail server for dev
  passreset/       — password reset CLI tool
internal/
  router/          — HTTP handlers and route registration
  auth/            — session/user auth helpers
  render/          — templ components
  mail/            — email sending
  readability/     — readability service client
  middlewares/     — auth protection, logging, etc.
  preview/         — link preview fetching
db/
  migrations/      — SQL migrations
  queries/         — sqlc-generated query code
  query.*.sql      — raw SQL queries
  schema.sql       — DB schema
readability/       — Python readability sidecar service
terraform/         — infrastructure (likely Hetzner + AWS)
```

## Running Locally

```bash
# Start services
docker-compose -f docker-compose.dev.yml up

# Or run Go server directly
go run ./cmd/pagemail
```

Config via `pagemail.yaml` or `PM_*` env vars. SQLite DB at `./db/pagemail.sqlite3`.

## Key Routes

| Route | Handler |
|-------|---------|
| `GET /pages` | List saved pages |
| `POST /pages` | Save a new page |
| `DELETE /pages/{id}` | Delete a page |
| `GET /pages/dashboard` | Main dashboard |
| `GET /user/account` | Account settings |
| `POST /login/google` | Google OAuth |
| `GET/POST /password-reset/*` | Password reset flow |

## Notes for Timmy

- The `articles`/`readings` tables were added and then removed. If Ellis wants to revisit that feature, the history is in PRs #157 and #160.
- The readability sidecar (`readability/`) is a separate Python service — changes there need a separate Docker image build.
- `sqlc.yaml` config means any DB schema changes require running `sqlc generate` to update the typed query code.
- Deployment is AWS ECR → Hetzner. Check `docker-compose.prd.yml` and `terraform/` for infra details.

# Linkrr

A fast, modular URL shortener and analytics service built with Go. Linkrr provides short URL creation, redirects, user management, and click analytics with a clean REST API.

---

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
  - [Option A: Clone or Fork](#option-a-clone-or-fork)
  - [Option B: Docker Pull](#option-b-docker-pull)
- [Configuration](#configuration)
- [Running](#running)
  - [Run Locally (Go)](#run-locally-go)
  - [Run with Docker](#run-with-docker)
- [API Documentation](#api-documentation)
- [Project Layout](#project-layout)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## Overview
Linkrr exposes REST endpoints for authentication, URL shortening, redirects, and analytics aggregation. It includes workers for email notifications and templates for login, signup, and password flows.

## Features
- Short URL creation with custom aliases
- Secure auth: signup, login, refresh tokens, password reset
- Redirect service to original URLs
- Analytics: click tracking, aggregation, global and per-link stats
- Email notifications with templated content
- Docker support for quick deployment

## Tech Stack
- Go (standard library + sqlc-generated data access)
- SQL schema and queries managed under `sql/`
- Caddy for reverse proxy (via `Caddyfile`)
- Docker for containerized deployment

## Quick Start

### Option A: Clone or Fork
Clone or fork the GitHub repository:

```bash
# Clone
git clone https://github.com/john-ayodeji/linkrr.git
cd linkrr

# Or fork on GitHub first, then clone your fork
```

### Option B: Docker Pull
Run the prebuilt image directly:

```bash
docker pull ayodejijohndev/linkrr:0.1.3
docker run --rm -p 8080:8080 --name linkrr ayodejijohndev/linkrr:0.1.3
```

Adjust ports as needed for your environment and proxy setup.

#### Enter container shell and prepare orchestration script
After pulling the image, you can enter the running container shell, copy the example orchestrator script, and run it to start required services:

```bash
# Enter the container shell
docker exec -it linkrr /bin/sh

# Inside the container, go to /app and copy the example script
cd /app
cp linkrr.example.sh linkrr.sh

# Edit linkrr.sh to set real secrets and credentials
# Then run it to start Postgres, one Linkrr instance, and Caddy
sh linkrr.sh
```

Notes:
- The `linkrr.example.sh` script starts one Linkrr container for testing.
- Caddy is used for load balancing and can accept up to three app instances (linkrr1–linkrr3). Edit `Caddyfile` in the container to add more upstream servers.

## Configuration
Key configuration files:
- `Caddyfile`: Reverse proxy configuration.
- `Dockerfile`: Container build instructions.
- `internal/config/apiConfig.go`: Application configuration (environment variables and defaults).

Environment variables you may want to set (examples; confirm against your deployment needs):
- `PORT`: Server listen port (e.g., `8080`).
- `DATABASE_URL`: Connection string for your database.
- `JWT_SECRET`: Secret for signing access tokens.
- `REFRESH_TOKEN_SECRET`: Secret for refresh tokens.
- `EMAIL_SMTP_HOST`, `EMAIL_SMTP_PORT`, `EMAIL_USERNAME`, `EMAIL_PASSWORD`: SMTP settings for outgoing emails.

## Running

### Run Locally (Go)
Ensure Go is installed and your database is reachable.

```bash
# From repo root
go mod download

# Start the server
go run ./...
```

You can also run the main entrypoint specifically:

```bash
go run ./main.go
```

### Run with Docker
Recommended approach using the prebuilt image and running the orchestrator on your host/WSL:

```bash
# 1) Pull the prebuilt image
docker pull ayodejijohndev/linkrr:0.1.3

# 2) Download the example orchestrator script (from GitHub)
curl -fsSL https://raw.githubusercontent.com/john-ayodeji/linkrr/refs/heads/main/linkrr.example.sh -o linkrr.sh
# Or (PowerShell)
# Invoke-WebRequest -Uri https://raw.githubusercontent.com/john-ayodeji/linkrr/refs/heads/main/linkrr.example.sh -OutFile linkrr.sh

# 3) Edit linkrr.sh to set real secrets (DB, JWT, SMTP, etc.)

# 4) Run the orchestrator (creates network, Postgres, one Linkrr app, and Caddy)
sh linkrr.sh
```

Notes:
- The example starts one app instance (`linkrr1`) for testing.
- Caddy handles load balancing and can be configured to accept up to three instances (`linkrr1`–`linkrr3`). Edit your `Caddyfile` to add upstreams.
- For production, run the script on your host with proper volumes, secrets management, monitoring, and backups.

## API Documentation
See [API_README.md](API_README.md) for endpoint details, request/response formats, and examples.

## Project Layout
A simplified overview of the structure:
- `internal/`
  - `auth/`: JWT & refresh token logic
  - `config/`: API configuration
  - `database/`: sqlc-generated DB access and models
  - `email_templates/`: HTML templates for emails
  - `events_workers/`: background workers (e.g., email sender)
  - `handlers/`: HTTP handlers grouped by domain
  - `services/`: business logic for auth, shortener, redirect, analytics, email, users
  - `utils/`: helpers for hashing, token parsing, error handling, template rendering
- `sql/`
  - `queries/`: SQL query definitions
  - `schema/`: migration files
- `main.go`: application entrypoint
- `Routes.go`: router setup
- `Dockerfile`: container build
- `Caddyfile`: reverse proxy config

## Development
Common tasks:

```bash
# Format
gofmt -s -w .

# Build
go build ./...

# Run
go run ./main.go
```

If you use sqlc, ensure it is installed and regenerate code after query/schema changes:

```bash
sqlc generate
```

Configuration and credentials should be provided via environment variables or `.env` files (avoid committing secrets).

## Testing
Run unit tests:

```bash
go test ./...
```

## Contributing
- Fork the repo and create a feature branch.
- Make changes with clear commit messages.
- Add/adjust tests where applicable.
- Open a pull request with a concise description of changes and rationale.



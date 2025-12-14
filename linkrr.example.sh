#!/bin/bash
set -euo pipefail

# Linkrr example orchestrator (sanitized)
# - Creates a docker network
# - Starts Postgres
# - Starts one Linkrr app container for testing
# - Starts Caddy as load balancer (can scale to 3 app instances)
#
# Copy to linkrr.sh and set real secrets before running.
#
# Required edits before production use:
# - Replace placeholder credentials and secrets below
# - Optionally scale app instances to 2-3 and update Caddyfile upstreams

NETWORK_NAME="linkrr-net"
POSTGRES_CONTAINER="my-postgres"
POSTGRES_IMAGE="postgres:15.5" # Use stable tag; adjust as needed
APP_IMAGE="ayodejijohndev/linkrr:0.1.3" # Prebuilt image from Docker Hub
APP_CONTAINER="linkrr1"
CADDY_CONTAINER="load-balancer"

# Placeholder credentials (CHANGE THESE)
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="CHANGE_ME_PASSWORD"
POSTGRES_DB="linkrr"
DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_CONTAINER:5432/$POSTGRES_DB?sslmode=disable"
JWT_SECRET="CHANGE_ME_JWT_SECRET"
REFRESH_TOKEN_SECRET="CHANGE_ME_REFRESH_SECRET"
MAILTRAP_TOKEN="CHANGE_ME_MAILTRAP"
IPSTACK_API_KEY="CHANGE_ME_IPSTACK"
INSTANCE_ID="001"

# Create network if not exists
docker network create "$NETWORK_NAME" 2>/dev/null || true

# Start Postgres
docker run -d \
  --name "$POSTGRES_CONTAINER" \
  --network "$NETWORK_NAME" \
  -e POSTGRES_USER="$POSTGRES_USER" \
  -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  -e POSTGRES_DB="$POSTGRES_DB" \
  -v pgdata:/var/lib/postgresql/data \
  -p 5432:5432 \
  "$POSTGRES_IMAGE"

# Give DB time to start
sleep 10

# Start single Linkrr app container for testing
docker run -d \
  --name "$APP_CONTAINER" \
  --network "$NETWORK_NAME" \
  -e DB_URL="$DB_URL" \
  -e JWT_SECRET="$JWT_SECRET" \
  -e REFRESH_TOKEN_SECRET="$REFRESH_TOKEN_SECRET" \
  -e MAILTRAP_TOKEN="$MAILTRAP_TOKEN" \
  -e IPSTACK_API_KEY="$IPSTACK_API_KEY" \
  -e INSTANCE_ID="$INSTANCE_ID" \
  --restart unless-stopped \
  "$APP_IMAGE"

# Start Caddy load balancer (reverse proxy)
# Caddyfile should define upstreams (accept up to 3)
# Edit the Caddyfile to add linkrr2, linkrr3 if you scale

docker run -d \
  --name "$CADDY_CONTAINER" \
  --network "$NETWORK_NAME" \
  -p 8080:80 \
  -v "$PWD/Caddyfile:/etc/caddy/Caddyfile" \
  caddy

echo "\nLinkrr stack started. Access via http://localhost:8080"
echo "To scale, start linkrr2/linkrr3 and update Caddyfile upstreams."
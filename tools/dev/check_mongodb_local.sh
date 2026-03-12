#!/usr/bin/env bash
set -euo pipefail

container_name="${ATHENA_MONGODB_CONTAINER:-mongodb-local}"
mongodb_uri="${ATHENA_MONGODB_URI:-mongodb://127.0.0.1:27017}"
mongodb_database="${ATHENA_MONGODB_DATABASE:-athenamind}"

if ! command -v podman >/dev/null 2>&1; then
  echo "podman is required" >&2
  exit 1
fi

if ! podman container exists "$container_name"; then
  echo "container not found: $container_name" >&2
  exit 1
fi

health_status="$(podman inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$container_name")"
container_status="$(podman inspect --format '{{.State.Status}}' "$container_name")"

if [[ "$container_status" != "running" ]]; then
  echo "container is not running: $container_name ($container_status)" >&2
  exit 1
fi

echo "mongodb_container: $container_name"
echo "container_status: $container_status"
echo "health_status: $health_status"
echo "mongodb_uri: $mongodb_uri"
echo "mongodb_database: $mongodb_database"

if podman exec "$container_name" mongosh --quiet --eval 'db.runCommand({ ping: 1 }).ok' >/tmp/athena-mongodb-ping.txt 2>/dev/null; then
  ping_result="$(tr -d '\r' </tmp/athena-mongodb-ping.txt | tail -n 1 | xargs)"
  echo "mongosh_ping_ok: $ping_result"
else
  echo "mongosh_ping_ok: unavailable"
  echo "note: container is running, but in-container mongosh ping did not succeed"
fi

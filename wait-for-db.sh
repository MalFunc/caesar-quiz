#!/bin/sh
set -e

# Tunggu Postgres siap. Host/port dari env DB_HOST/DB_PORT.
host="${DB_HOST:-db}"
port="${DB_PORT:-5432}"

until pg_isready -h "$host" -p "$port" >/dev/null 2>&1; do
  >&2 echo "Postgres ($host:$port) belum siap - sleeping"
  sleep 1
done

>&2 echo "Postgres siap - menjalankan: $*"
exec "$@"

#!/bin/bash
set -e

host="$DB_HOST"
port="$DB_PORT"
user="$DB_USER"
password="$DB_PASSWORD"
db="$DB_NAME"
connection_string="postgres://$user:$password@$host:$port/$db?sslmode=disable"

echo "Using connection string: $connection_string"

attempt=1
max_attempts=30

until psql "$connection_string" -c '\q' 2>/dev/null; do
  echo "Attempt $attempt/$max_attempts: Waiting for PostgreSQL to be ready at $host:$port..."
  if [ $attempt -ge $max_attempts ]; then
    echo "Error: PostgreSQL is not ready after $max_attempts attempts"
    exit 1
  fi
  sleep 2
  ((attempt++))
done

echo "PostgreSQL is ready!"
exec "$@"
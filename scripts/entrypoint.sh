#!/bin/sh
set -e

# Wait for postgres to be ready
echo "Waiting for postgres..."
until nc -z postgres 5432; do
  sleep 1
done
echo "Postgres is up!"

# Run migrations (if we had a migration tool in the binary or separate)
# For now, let's assume we can run it via go run if env allows, 
# but since this is production-ready, let's check for a migrate command.
# Since the project has cmd/migrate, we could build that too and run it.

echo "Starting application..."
exec "$@"

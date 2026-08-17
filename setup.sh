#!/bin/bash
set -e

echo "🚀 Starting FinApp Full Stack Setup..."

# 1. Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Error: docker-compose not found. Please install Docker and Docker Compose."
    exit 1
fi

# 2. Pull latest images for infrastructure services
echo "📥 Pulling infrastructure images (Postgres, Redis, Minio)..."
docker compose pull postgres minio redis

# 3. Build and Start Services
echo "🏗️ Building and starting all services..."
docker compose up -d --build

echo "⏳ Waiting for services to be healthy..."
# Wait for the api service specifically as it depends on postgres healthcheck
timeout 60s bash -c 'until [ "$(docker inspect -f {{.State.Health.Status}} finapp-api)" == "healthy" ]; do sleep 2; done' || {
  echo "❌ Error: API failed to become healthy within 60 seconds."
  docker compose logs api
  exit 1
}

echo ""
echo "✅ FinApp is up and running!"
echo "--------------------------------------------------"
echo "🌐 Frontend: http://localhost:3000"
echo "⚙️  Backend API: http://localhost:8080"
echo "📦 MinIO Console: http://localhost:9001"
echo "<0xF0><0x9F><0x97><0x83>️  Redis: localhost:6379"
echo "🐘 Postgres: localhost:5432"
echo "--------------------------------------------------"
echo ""

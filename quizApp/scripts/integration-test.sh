#!/bin/bash
set -e

echo "Starting integration tests..."

# Start services
echo "Starting Docker containers..."
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be ready
echo "Waiting for API to be ready..."
max_attempts=30
attempt=0
until curl -s http://localhost:8080/categories/ > /dev/null || [ $attempt -eq $max_attempts ]; do
  echo "Waiting for API... (attempt $((attempt+1))/$max_attempts)"
  sleep 2
  attempt=$((attempt+1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "API failed to start"
  docker-compose -f docker-compose.test.yml logs
  docker-compose -f docker-compose.test.yml down
  exit 1
fi

# Run integration tests
echo "Running integration tests..."
go test -tags=integration ./api/... -v

# Capture exit code
TEST_EXIT_CODE=$?

# Cleanup
echo "Cleaning up..."
docker-compose -f docker-compose.test.yml down -v

exit $TEST_EXIT_CODE

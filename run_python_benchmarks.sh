#!/bin/bash

echo "Running Python RabbitMQ benchmarks..."

# Start RabbitMQ
echo "Starting RabbitMQ..."
docker compose up -d

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to be ready..."
sleep 10

# Run Python benchmarks
echo "Running benchmarks..."
docker compose -f docker/docker-compose-all.yml run --build --rm producer-fast

echo "Benchmarks complete!"

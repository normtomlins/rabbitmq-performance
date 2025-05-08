#!/bin/bash

echo "Running Go RabbitMQ benchmarks..."

# Start RabbitMQ
echo "Starting RabbitMQ..."
docker compose up -d

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to be ready..."
sleep 10

# Run Go benchmarks
echo "Running benchmarks..."
docker compose -f docker/docker-compose-go-simple.yml run --build --rm go-producer

echo "Benchmarks complete!"

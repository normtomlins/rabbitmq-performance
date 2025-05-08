#!/bin/bash

echo "RabbitMQ Go Producer Test (Simple)"
echo "================================="
echo ""

# Make sure RabbitMQ is running
docker compose -f docker-compose-go-simple.yml up -d rabbitmq

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to be ready..."
sleep 5

if [ "$1" == "standard" ]; then
    echo "Running standard Go producer..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer
elif [ "$1" == "ultra" ]; then
    echo "Running ultra-fast Go producer..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run ultra_fast.go
else
    echo "Running standard Go producer (default)..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer
fi

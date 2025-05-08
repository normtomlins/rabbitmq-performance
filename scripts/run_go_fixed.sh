#!/bin/bash

echo "RabbitMQ Go Producer Test Suite"
echo "==============================="
echo ""

# Make sure RabbitMQ is running
docker compose -f docker-compose-go-simple.yml up -d rabbitmq

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to be ready..."
sleep 5

if [ "$1" == "standard" ]; then
    echo "Running standard Go producer..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run main.go
elif [ "$1" == "ultra-fixed" ]; then
    echo "Running fixed ultra-fast Go producer..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run ultra_fast_fixed.go
elif [ "$1" == "batch" ]; then
    echo "Running ultra batch Go producer..."
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run ultra_batch.go
elif [ "$1" == "all" ]; then
    echo "Running all Go producer tests..."
    echo ""
    
    echo "=== TEST 1: Standard Go Producer ==="
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run main.go
    
    echo ""
    echo "Pausing between tests..."
    sleep 3
    echo ""
    
    echo "=== TEST 2: Ultra Fast Go Producer (Fixed) ==="
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run ultra_fast_fixed.go
    
    echo ""
    echo "Pausing between tests..."
    sleep 3
    echo ""
    
    echo "=== TEST 3: Ultra Batch Go Producer ==="
    docker compose -f docker-compose-go-simple.yml run --build --rm go-producer go run ultra_batch.go
    
    echo ""
    echo "=== All tests completed! ==="
else
    echo "Usage: $0 [standard|ultra-fixed|batch|all]"
    echo ""
    echo "Options:"
    echo "  standard    - Run the standard Go producer with multiple tests"
    echo "  ultra-fixed - Run the fixed ultra-fast Go producer"
    echo "  batch       - Run the ultra batch Go producer"
    echo "  all         - Run all tests in sequence"
fi

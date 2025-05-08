#!/bin/bash

echo "RabbitMQ Concurrent Producer Tests"
echo "================================="
echo ""

if [ "$1" == "threaded" ]; then
    echo "Running THREADED producer..."
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast python producer_threaded_docker.py
elif [ "$1" == "async" ]; then
    echo "Running ASYNC producer (requires aio-pika)..."
    # First install aio-pika
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast pip install aio-pika
    # Then run the async producer
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast python producer_async_docker.py
elif [ "$1" == "concurrent" ]; then
    echo "Running CONCURRENT producer (ThreadPoolExecutor)..."
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast python producer_concurrent_docker.py
elif [ "$1" == "test-all" ]; then
    echo "Testing all concurrent approaches..."
    
    echo -e "\n1. THREADED PRODUCER TEST"
    echo "========================="
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast python producer_threaded_docker.py test
    
    echo -e "\n2. CONCURRENT PRODUCER TEST"
    echo "==========================="
    docker compose -f docker-compose-fast.yml run --build --rm -v $(pwd):/app producer-fast python producer_concurrent_docker.py test
else
    echo "Usage: $0 [threaded|async|concurrent|test-all]"
    echo ""
    echo "Options:"
    echo "  threaded    - Run with Python threading (like Go goroutines)"
    echo "  async       - Run with Python asyncio (requires aio-pika)"
    echo "  concurrent  - Run with ThreadPoolExecutor"
    echo "  test-all    - Test all approaches with different configurations"
fi

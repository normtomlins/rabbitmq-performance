#!/bin/bash

echo "Initializing Go module..."
cd /Users/norm/Desktop/rabbitmq-python-demo/go-producer

# Initialize module if needed
if [ ! -f "go.mod" ]; then
    go mod init rabbitmq-go-producer
fi

# Download dependencies
go mod download
go mod tidy

echo "Go module initialized successfully!"

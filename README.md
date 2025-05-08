# RabbitMQ Performance: Python vs Go

A comprehensive benchmark comparing RabbitMQ message throughput between Python and Go implementations.

## Overview

This project demonstrates the performance differences between Python and Go when publishing messages to RabbitMQ. Through various optimizations and concurrency models, we discovered dramatic performance differences between the two languages.

### Key Findings

- **Python's best**: 72,204 messages/second (multiprocessing)
- **Go's peak**: 1,082,014 messages/second (ultra-optimized)
- **Go's sustained**: 266,025 messages/second
- Python's GIL severely limits threading performance
- Go's goroutines provide true concurrency

## Project Structure

```
.
├── python/                 # Python implementations
│   ├── producer.py        # Basic producer
│   ├── consumer.py        # Basic consumer
│   ├── producer_*.py      # Various optimizations
│   └── requirements.txt
├── go-producer/           # Go implementations
│   ├── main.go           # Standard tests
│   ├── ultra_fast.go     # Optimized version
│   ├── go.mod
│   └── go.sum
├── docker-compose*.yml    # Docker configurations
├── blog_post.html        # Main blog post
├── blog_technical_appendix.html  # Technical details
└── README.md
```

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Python 3.9+ (optional, for local development)
- Go 1.21+ (optional, for local development)

### Running the Benchmarks

1. Clone the repository:

   ```bash
   git clone https://github.com/normtomlins/rabbitmq-performance.git
   cd rabbitmq-performance
   ```

2. Start RabbitMQ:

   ```bash
   docker compose up -d rabbitmq
   ```

3. Run Python benchmarks:

   ```bash
   ./run_go_fixed.sh all
   ```

4. Run Go benchmarks:
   ```bash
   ./run_go_fixed.sh all
   ```

## Performance Results

### Python Performance

| Implementation                | Messages/Second | Notes                     |
| ----------------------------- | --------------- | ------------------------- |
| Single Thread (Persistent)    | ~1,400          | With publisher confirms   |
| Single Thread (Optimized)     | ~22,000         | Non-persistent            |
| Multi-threading (8 threads)   | ~5,900          | Worse than single thread! |
| Multiprocessing (4 processes) | **72,204**      | Best Python performance   |

### Go Performance

| Implementation              | Messages/Second | Notes                      |
| --------------------------- | --------------- | -------------------------- |
| Single Goroutine            | 220,773         | 3x faster than best Python |
| Worker Pool (10 workers)    | 81,210          | Connection contention      |
| Ultra-optimized (peak)      | **1,082,014**   | 15x faster than Python     |
| Ultra-optimized (sustained) | 266,025         | Average over 500K messages |

## Key Optimizations

### Python Optimizations

- Remove console output (+20x performance)
- Use multiprocessing instead of threading
- Non-persistent messages
- Connection per process

### Go Optimizations

- Connection pooling
- Pre-allocation of objects
- Simplified message structure
- Optimal worker count (2x CPU cores)
- Non-persistent messages

## Blog Posts

- [Main Blog Post](blog_post.html) - Comprehensive overview of findings
- [Technical Appendix](blog_technical_appendix.html) - Detailed code analysis

## Contributing

Feel free to run these benchmarks on your own hardware and share your results!

## License

MIT License - see LICENSE file for details

## Acknowledgments

- RabbitMQ team for excellent documentation
- Python `pika` library maintainers
- Go `amqp091-go` library maintainers

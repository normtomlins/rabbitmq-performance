# RabbitMQ Performance Benchmark: Python vs Go

A comprehensive benchmark comparing RabbitMQ message throughput between Python and Go implementations, revealing dramatic performance differences and optimization insights.

## 🚀 Quick Results

- **Python's Best**: 72,204 messages/second (multiprocessing)
- **Go's Peak**: 1,082,014 messages/second (15x faster!)
- **Go's Sustained**: 266,025 messages/second
- **Key Finding**: Python's GIL severely limits threading performance

## 📊 Blog Posts

- [**Main Analysis**](docs/blog_post.html) - Comprehensive overview with charts and findings
- [**Technical Deep Dive**](docs/blog_technical_appendix.html) - Code analysis and optimization details

## 🏃‍♂️ Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/rabbitmq-performance.git
cd rabbitmq-performance

# Start RabbitMQ
docker compose up -d

# Run Python benchmarks
./run_python_benchmarks.sh

# Run Go benchmarks
./run_go_benchmarks.sh
```

## 📁 Project Structure

```
.
├── python/              # Python implementations
│   ├── producer.py     # Basic producer
│   ├── consumer.py     # Basic consumer
│   └── ...            # Various optimizations
├── go/                # Go implementations
│   ├── main.go        # Standard benchmarks
│   ├── ultra_fast.go  # Optimized version
│   └── ...
├── docker/            # Docker configurations
├── docs/              # Documentation & blog posts
├── scripts/           # Helper scripts
└── index.html         # Project landing page
```

## 🔍 Key Findings

1. **Python's GIL is a major bottleneck** - Multi-threading actually decreased performance
2. **Multiprocessing saves Python** - Achieved 72K msgs/sec by bypassing the GIL
3. **Go's concurrency shines** - Goroutines enable true parallelism
4. **More threads ≠ better performance** - Optimization matters more than thread count
5. **Small changes, big impact** - Removing console output improved performance 20x

## 📈 Performance Comparison

### Python Results

| Implementation | Messages/Second | Notes |
|---------------|-----------------|-------|
| Single Thread (Optimized) | ~22,000 | Non-persistent messages |
| Multi-threading (8 threads) | ~5,900 | **Worse** than single thread! |
| Multiprocessing (4 processes) | **72,204** | Best Python performance |

### Go Results

| Implementation | Messages/Second | Notes |
|---------------|-----------------|-------|
| Single Goroutine | 220,773 | 3x faster than best Python |
| Ultra-optimized (peak) | **1,082,014** | 15x faster than Python |
| Ultra-optimized (sustained) | 266,025 | Average over 500K messages |

## 🛠 Technologies Used

- Python 3.9 with `pika` library
- Go 1.21 with `amqp091-go` library  
- RabbitMQ 3.x
- Docker & Docker Compose
- macOS with Apple Silicon (M1)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Run the benchmarks on your hardware and share your results!

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- RabbitMQ team for excellent documentation
- Python `pika` library maintainers
- Go `amqp091-go` library maintainers
- The developer community for feedback and suggestions

---

**Questions?** Open an issue or reach out!

**Found this useful?** Give it a ⭐️!

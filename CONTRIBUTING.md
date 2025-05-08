# Contributing to RabbitMQ Performance Benchmark

We welcome contributions to this benchmarking project! Here's how you can help:

## Ways to Contribute

1. **Run benchmarks on different hardware**
   - Test on different CPU architectures (ARM, x86)
   - Try different operating systems
   - Test with various RabbitMQ configurations

2. **Add new implementations**
   - Try other languages (Rust, Java, C++)
   - Implement different messaging patterns
   - Test with RabbitMQ clusters

3. **Improve existing code**
   - Optimize current implementations
   - Fix bugs or issues
   - Improve documentation

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-optimization`)
3. Make your changes
4. Run the benchmarks to ensure everything works
5. Commit your changes (`git commit -m 'Add amazing optimization'`)
6. Push to your branch (`git push origin feature/amazing-optimization`)
7. Create a Pull Request

## Benchmark Guidelines

When adding new benchmarks:

1. Use consistent message formats
2. Test with both persistent and non-persistent messages
3. Document your hardware and software environment
4. Run tests multiple times to ensure consistency
5. Include both peak and sustained performance metrics

## Code Style

### Python
- Follow PEP 8
- Use type hints where appropriate
- Include docstrings for functions

### Go
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names
- Add comments for complex logic

## Reporting Results

When sharing benchmark results, please include:

- Hardware specifications (CPU, RAM, Storage type)
- Operating system and version
- Docker version
- RabbitMQ version
- Network configuration (if relevant)
- Any special settings or optimizations

## Questions?

Feel free to open an issue for any questions or discussions!

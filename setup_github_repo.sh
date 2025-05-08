#!/bin/bash

echo "Setting up GitHub repository..."

# Initialize git
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: RabbitMQ Performance Benchmark - Python vs Go"

# Create main branch
git branch -M main

echo ""
echo "Repository initialized!"
echo ""
echo "Next steps:"
echo "1. Create a new repository on GitHub named 'rabbitmq-performance'"
echo "2. Run the following commands:"
echo ""
echo "   git remote add origin https://github.com/normtomlins/rabbitmq-performance.git"
echo "   git push -u origin main"
echo ""
echo "3. Enable GitHub Pages:"
echo "   - Go to repository Settings > Pages"
echo "   - Select 'Deploy from a branch'"
echo "   - Choose 'main' branch and '/ (root)' folder"
echo "   - Save"
echo ""
echo "Your benchmark results will be available at:"
echo "https://normtomlins.github.io/rabbitmq-performance/"

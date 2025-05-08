#!/bin/bash

echo "Removing workflow file temporarily..."

# Remove the workflow directory
rm -rf .github

# Add and commit the removal
git add .
git commit -m "Remove workflow file temporarily"

# Push to GitHub
echo "Pushing to GitHub..."
git push -u origin main

echo ""
echo "Push complete! You can add the workflow file later through GitHub's web interface."
echo ""
echo "To add workflow later:"
echo "1. Go to your repository on GitHub"
echo "2. Click 'Actions' tab"
echo "3. Click 'New workflow'"
echo "4. Choose 'set up a workflow yourself'"
echo "5. Copy the content from the workflow file"

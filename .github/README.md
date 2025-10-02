# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automating various tasks in the repository.

## Workflows

### Publish Docker Image
- **Trigger**: Runs on pushes to the `main` branch and when a release is published
- **Purpose**: Builds and publishes the Docker image to GitHub Container Registry (GHCR)
- **Platforms**: Builds for both AMD64 and ARM64 architectures
- **Tags**: Creates multiple tags including branch, SHA, semantic versions, and latest

## Registry

The Docker image is published to: `ghcr.io/username/city-autocomplete-api`

## Configuration

No additional configuration is required. The workflow uses the default GITHUB_TOKEN for authentication.
# Deploying on Coolify

This guide explains how to deploy the City Autocomplete API on Coolify, a self-hosted alternative to Vercel/Netlify.

## Prerequisites

- Access to a Coolify instance
- This repository connected to your Coolify installation

## Deployment Steps

1. **Add the Application**:
   - Navigate to your Coolify dashboard
   - Click "Add Application"
   - Select your Git provider and repository (frostylabs/city-autocomplete-api)

2. **Configure Build Settings**:
   - Build Pack: `Dockerfile`
   - Build Path: `.` (root of repository)
   - Dockerfile Path: `Dockerfile` (default)
   - Build Context: `.` (root of repository)

3. **Configure Environment Variables** (optional):
   - `PORT`: 8080 (or desired port)
   - `DB_PATH`: cities.db
   - `CSV_PATH`: world-cities.csv

4. **Set Up Port Mapping**:
   - Map internal port `8080` to your desired external port

5. **Deploy**:
   - Click "Save & Deploy"

## Database Persistence

The application will create and populate the database on first startup. For data persistence:

1. Set up a volume mount for the database file
2. Mount to path: `/app/cities.db`

## Health Check

Set the health check path to:
- `/` (root endpoint returns basic service info)

## Environment-Specific Configurations

You can modify these environment variables to customize your deployment:

- `PORT`: Change from default 8080 if needed
- `DB_PATH`: Path to SQLite database (default: cities.db)
- `CSV_PATH`: Path to CSV file for initial data (default: world-cities.csv)

## Notes

- The first deployment may take longer as it needs to build the database from the CSV file
- Subsequent deployments will be faster as the database is already built
- The Docker image includes the pre-populated database to speed up initial startup
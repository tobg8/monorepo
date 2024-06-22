#!/bin/bash

# Pull the latest PostgreSQL Docker image
docker pull postgres

# Create a Docker volume for PostgreSQL data
docker volume create postgres_data

# Run PostgreSQL container
docker run --name postgres_container \
           -e POSTGRES_PASSWORD=mysecretpassword \
           -d -p 5432:5432 \
           -v postgres_data:/var/lib/postgresql/data \
           postgres

echo "PostgreSQL Docker container is now running."
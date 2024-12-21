#!/bin/sh

# Define default uid and gid to match the host user
BUILD_UID=1027
BUILD_GID=100

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo "Docker is not running. Starting Docker..."
        if [ "$(uname)" = "Darwin" ]; then
            # macOS
            open --background -a Docker
        elif [ "$(uname)" = "Linux" ]; then
            # Linux
            sudo systemctl start docker
        elif [ "$(uname)" = "MINGW64_NT-10.0" ]; then
            # Windows
            start-service docker
        else
            echo "Unsupported OS. Please start Docker manually."
            exit 1
        fi

        # Wait for Docker to start
        while ! docker info > /dev/null 2>&1; do
            echo "Waiting for Docker to start..."
            sleep 2
        done
        echo "Docker started."
    else
        echo "Docker is already running."
    fi
}

# Check if Docker is running
check_docker

# Build the Docker image for the specified platform with uid and gid as build arguments
docker buildx build --platform linux/amd64 --build-arg UID=${BUILD_UID} --build-arg GID=${BUILD_GID} -t photo-server .

# Save the Docker image to a tar file
docker save photo-server -o photo-server.tar

echo "Docker image built and saved to photo-server.tar"

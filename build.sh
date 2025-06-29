#!/bin/bash

# Build script for pharmacy application with latest code pull

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="pharmacy-app"
TAG="latest"
GIT_REPO_URL=""  # Set your git repository URL here
GIT_BRANCH="main"
GIT_TOKEN=""     # Set your git token if using private repo

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -r, --repo URL        Git repository URL"
    echo "  -b, --branch BRANCH   Git branch (default: main)"
    echo "  -t, --token TOKEN     Git access token for private repos"
    echo "  -i, --image NAME      Docker image name (default: pharmacy-app)"
    echo "  --tag TAG             Docker image tag (default: latest)"
    echo "  --local               Build from local directory (default behavior)"
    echo "  --git                 Build from git repository"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --local                                    # Build from local directory"
    echo "  $0 --git -r https://github.com/user/repo.git # Build from git repo"
    echo "  $0 --git -r https://github.com/user/repo.git -b develop -t ghp_token"
}

# Parse command line arguments
BUILD_FROM_GIT=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -r|--repo)
            GIT_REPO_URL="$2"
            shift 2
            ;;
        -b|--branch)
            GIT_BRANCH="$2"
            shift 2
            ;;
        -t|--token)
            GIT_TOKEN="$2"
            shift 2
            ;;
        -i|--image)
            IMAGE_NAME="$2"
            shift 2
            ;;
        --tag)
            TAG="$2"
            shift 2
            ;;
        --local)
            BUILD_FROM_GIT=false
            shift
            ;;
        --git)
            BUILD_FROM_GIT=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

print_status "Starting Docker build process..."

if [ "$BUILD_FROM_GIT" = true ]; then
    if [ -z "$GIT_REPO_URL" ]; then
        print_error "Git repository URL is required when building from git"
        show_usage
        exit 1
    fi
    
    print_status "Building from git repository: $GIT_REPO_URL"
    print_status "Branch: $GIT_BRANCH"
    
    # Build arguments for Docker
    BUILD_ARGS="--build-arg GIT_REPO_URL=$GIT_REPO_URL --build-arg GIT_BRANCH=$GIT_BRANCH"
    
    if [ -n "$GIT_TOKEN" ]; then
        BUILD_ARGS="$BUILD_ARGS --build-arg GIT_TOKEN=$GIT_TOKEN"
        print_status "Using authentication token"
    fi
    
    # Add git configuration from global settings if available
    GIT_USER_NAME=$(git config --global user.name 2>/dev/null || echo "")
    GIT_USER_EMAIL=$(git config --global user.email 2>/dev/null || echo "")
    
    if [ -n "$GIT_USER_NAME" ]; then
        BUILD_ARGS="$BUILD_ARGS --build-arg GIT_USER_NAME='$GIT_USER_NAME'"
        print_status "Using git user name: $GIT_USER_NAME"
    fi
    
    if [ -n "$GIT_USER_EMAIL" ]; then
        BUILD_ARGS="$BUILD_ARGS --build-arg GIT_USER_EMAIL='$GIT_USER_EMAIL'"
        print_status "Using git user email: $GIT_USER_EMAIL"
    fi
    
    # Build using the git Dockerfile
    print_status "Building Docker image: $IMAGE_NAME:$TAG"
    docker build -f Dockerfile.git $BUILD_ARGS -t "$IMAGE_NAME:$TAG" .
    
else
    print_status "Building from local directory"
    
    # Pull latest changes if in a git repository
    if [ -d ".git" ]; then
        print_status "Pulling latest changes from git..."
        git pull origin $(git branch --show-current) || print_warning "Could not pull latest changes"
    fi
    
    # Update go modules
    print_status "Updating Go modules..."
    go mod tidy
    go mod download
    
    # Add git configuration from global settings for local builds too
    BUILD_ARGS=""
    GIT_USER_NAME=$(git config --global user.name 2>/dev/null || echo "")
    GIT_USER_EMAIL=$(git config --global user.email 2>/dev/null || echo "")
    
    if [ -n "$GIT_USER_NAME" ]; then
        BUILD_ARGS="$BUILD_ARGS --build-arg GIT_USER_NAME='$GIT_USER_NAME'"
        print_status "Using git user name: $GIT_USER_NAME"
    fi
    
    if [ -n "$GIT_USER_EMAIL" ]; then
        BUILD_ARGS="$BUILD_ARGS --build-arg GIT_USER_EMAIL='$GIT_USER_EMAIL'"
        print_status "Using git user email: $GIT_USER_EMAIL"
    fi
    
    # Build using the standard Dockerfile
    print_status "Building Docker image: $IMAGE_NAME:$TAG"
    docker build $BUILD_ARGS -t "$IMAGE_NAME:$TAG" .
fi

if [ $? -eq 0 ]; then
    print_status "Docker build completed successfully!"
    print_status "Image: $IMAGE_NAME:$TAG"
    
    # Show image details
    docker images | grep "$IMAGE_NAME" | head -1
    
    print_status "To run the container:"
    echo "  docker run -p 8080:8080 $IMAGE_NAME:$TAG"
    echo "  or"
    echo "  docker-compose up"
else
    print_error "Docker build failed!"
    exit 1
fi 
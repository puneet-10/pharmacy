# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git (required for some Go modules and git operations)
RUN apk add --no-cache git

# Set up git configuration from build args (will use global config if available)
ARG GIT_USER_NAME
ARG GIT_USER_EMAIL
RUN if [ -n "$GIT_USER_NAME" ]; then git config --global user.name "$GIT_USER_NAME"; fi
RUN if [ -n "$GIT_USER_EMAIL" ]; then git config --global user.email "$GIT_USER_EMAIL"; fi

# If building from a git repository, clone the latest code
# Replace YOUR_GIT_REPO_URL with your actual repository URL
# ARG GIT_REPO_URL
# ARG GIT_BRANCH=main
# RUN git clone --depth 1 --branch ${GIT_BRANCH} ${GIT_REPO_URL} .

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Run go mod vendor to create vendor directory
RUN go mod vendor

# Copy the source code into the container
COPY . .

# Verify dependencies are up to date
RUN go mod tidy

# Build the application with vendor directory
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership of the app directory
RUN chown -R appuser:appgroup /root

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 
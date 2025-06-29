# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Install git and other necessary tools
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Build arguments for git repository
ARG GIT_REPO_URL
ARG GIT_BRANCH=main
ARG GIT_TOKEN
ARG GIT_USER_NAME
ARG GIT_USER_EMAIL

# Set up git configuration from build args (will use global config if available)
RUN if [ -n "$GIT_USER_NAME" ]; then git config --global user.name "$GIT_USER_NAME"; fi
RUN if [ -n "$GIT_USER_EMAIL" ]; then git config --global user.email "$GIT_USER_EMAIL"; fi

# Clone the repository with the latest code
RUN if [ -n "$GIT_TOKEN" ]; then \
        git clone --depth 1 --branch ${GIT_BRANCH} https://${GIT_TOKEN}@${GIT_REPO_URL#https://} . ; \
    else \
        git clone --depth 1 --branch ${GIT_BRANCH} ${GIT_REPO_URL} . ; \
    fi

# Download dependencies
RUN go mod download

# Create vendor directory
RUN go mod vendor

# Verify and tidy dependencies
RUN go mod tidy

# Build the application using vendor directory
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
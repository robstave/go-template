# === Stage 1: Build the React Frontend ===
FROM node:18-alpine AS frontend
WORKDIR /app
# Copy frontend package files and install dependencies
COPY meowmorize-frontend/package*.json ./meowmorize-frontend/
RUN cd meowmorize-frontend && npm install
# Copy all frontend source files and build
COPY meowmorize-frontend/ ./meowmorize-frontend/
RUN cd meowmorize-frontend && npm run build

# === Stage 2: Build the Go Backend ===
FROM golang:1.23-alpine AS builder
WORKDIR /app
# Enable cgo and set build parameters
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
# Install necessary packages for building (gcc, musl-dev, etc.)
RUN apk update && apk add --no-cache git gcc musl-dev
# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the entire project
COPY . .
# Copy the React build output from the frontend stage into the expected directory
COPY --from=frontend /app/meowmorize-frontend/build ./meowmorize-frontend/build
# Build the Go binary (ensure your main.go serves static files from ./meowmorize-frontend/build)
RUN go build -o backend ./cmd/main/main.go

# === Stage 3: Create the Final Production Image ===
FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
# Copy the Go binary from the builder stage
COPY --from=builder /app/backend .
# Copy the React build output
COPY --from=builder /app/meowmorize-frontend/build ./meowmorize-frontend/build
# Create a directory for SQLite data
RUN mkdir -p /app/data
# Expose the desired port (e.g., 8999)
EXPOSE 8999
# Set environment variables as needed
ENV DB_PATH=/app/data/db.sqlite3
ENV PORT=8999
# Add a non-root user for security
# Optional default values for the initial user credentials
# ENV DEFAULT_USER_USERNAME=meow
# ENV DEFAULT_USER_PASSWORD=meow
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && chown -R appuser:appgroup /app
USER appuser
# Run the backend binary
CMD ["./backend"]

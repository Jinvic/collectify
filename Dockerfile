# Use a multi-stage build to optimize the image size
# Stage 1: Build the React frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy package.json and pnpm-lock.yaml for dependency installation
# This step is cached if these files haven't changed
COPY web/package.json web/pnpm-lock.yaml ./web/

# Install pnpm globally and then install frontend dependencies
RUN npm install -g pnpm && \
    cd web && \
    pnpm install --frozen-lockfile

# Copy the rest of the frontend source code
COPY web/ ./web/

# Build the React application
RUN cd web && pnpm run build

# Stage 2: Build the Go backend
# Use the official Golang Alpine image for a smaller footprint
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Install Git (needed by Go modules)
# Note: We do NOT install gcc or musl-dev anymore because github.com/glebarez/sqlite is pure Go
RUN apk add --no-cache git

# Copy go mod and sum files for dependency download
# This step is cached if these files haven't changed
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary.
# -ldflags="-w -s" strips debug symbols to reduce binary size
# Explicitly set CGO_ENABLED=0 to ensure a pure Go build, leveraging the glebarez/sqlite driver's advantage.
# GOOS=linux ensures the binary is built for the Linux target OS inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o collectify .

# Stage 3: Final stage - Create the minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS requests (if needed by the app in the future)
RUN apk --no-cache add ca-certificates

# Create a non-root user for security
RUN adduser -D -s /bin/sh collectify-user

# Copy the pre-built binary file from the backend-builder stage
COPY --from=backend-builder /app/collectify .

# Copy the built frontend static files from the frontend-builder stage
# The Go backend is configured to serve files from ./web/build
COPY --from=frontend-builder /app/web/build ./web/build

# Make the binary executable
RUN chmod +x ./collectify

# Change ownership of the binary and static files to the non-root user
RUN chown -R collectify-user:collectify-user ./collectify ./web/build

# Switch to the non-root user
USER collectify-user

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./collectify"]
CMD ["web"]
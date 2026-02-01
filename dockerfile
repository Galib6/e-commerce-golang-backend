# ---- Build stage ----
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install git (needed for some go modules)
RUN apk add --no-cache git

# Copy go mod files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# ---- Production stage ----
FROM alpine:3.20

WORKDIR /app

# Certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Copy compiled binary
COPY --from=builder /app/app .

# Copy migrations folder for TypeORM-style migrations
COPY --from=builder /app/migrations ./migrations

# Expose app port (change if needed)
EXPOSE 3000

# Run the binary
CMD ["./app"]



#  docker build -t go-api .
#  docker run -p 8080:8080 go-api
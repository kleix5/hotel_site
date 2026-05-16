# Multi-stage build to keep final runtime image small and focused.
FROM golang:1.23-alpine AS build

# Build context root for compilation stage.
WORKDIR /app

# Copy dependency manifest first to maximize Docker layer cache reuse.
COPY go.mod ./
RUN go mod download

# Copy full source and build static Linux binary.
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# Minimal runtime image containing only compiled server and runtime assets.
FROM alpine:3.20

# Runtime working directory used by binary and template/static file paths.
WORKDIR /app

# Copy only artifacts needed at runtime: binary, templates/static assets, and schema files.
COPY --from=build /bin/server /app/server
COPY --from=build /app/web /app/web
COPY --from=build /app/db /app/db

# App listens on 8080 inside container.
EXPOSE 8080

# Start HTTP server process.
CMD ["/app/server"]

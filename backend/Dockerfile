# Multi-stage build for the wonderlic-calc Go backend.

# --- build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src

# Cache deps first.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- runtime stage ---
FROM alpine:3.20
# wget (busybox) is used by the compose healthcheck against /healthz.
WORKDIR /app
COPY --from=build /out/server /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]

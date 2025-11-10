# Stage 1: Build the application
FROM golang:1.25-alpine AS build

RUN apk update && \
    apk add build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o fuelbot ./cmd/fuelbot

# Stage 2: Create a minimal runtime image
FROM alpine:latest

RUN apk update && apk add --no-cache \
    ca-certificates  \
    sqlite

WORKDIR /app
COPY --from=build /app/fuelbot /usr/local/bin/fuelbot

VOLUME ["/var/lib/fuelbot"]

ENTRYPOINT ["/usr/local/bin/fuelbot"]

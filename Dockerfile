# Stage 1: Build the application with CGO and GCC
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
WORKDIR /app
VOLUME ["/var/lib/fuelbot"]
COPY --from=build /app/fuelbot /usr/local/bin/fuelbot
ENTRYPOINT ["/usr/local/bin/fuelbot"]

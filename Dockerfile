FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fuelbot ./cmd/fuelbot

FROM alpine:latest
WORKDIR /app
ENV DB_PATH=/data/fuelbot.db
VOLUME ["/data"]
COPY --from=build /app/fuelbot /usr/local/bin/fuelbot
ENTRYPOINT ["/usr/local/bin/fuelbot"]
# ожидается TELEGRAM_BOT_TOKEN через env

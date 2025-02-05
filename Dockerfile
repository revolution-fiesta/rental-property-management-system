FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o rental-app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rental-app .
COPY config/config.yaml ./config/
EXPOSE 8080
CMD ["./rental-app"]
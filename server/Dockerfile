# Stage 1 - build
FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o gochat cmd/main.go

# Stage 2 - run
FROM alpine
WORKDIR /app
COPY --from=builder /app/gochat .
EXPOSE 8080
CMD ["./gochat"]
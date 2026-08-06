# Build stage
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o echallan-calculator cmd/echallan-calculator/main.go

# Final stage
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/echallan-calculator .
COPY --from=builder /app/configs ./configs

EXPOSE 8078

CMD ["./echallan-calculator"]

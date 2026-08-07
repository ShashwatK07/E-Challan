# Build stage
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o echallan-services cmd/echallan-services/main.go

# Final stage
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/echallan-services .
COPY --from=builder /app/configs ./configs

EXPOSE 8079

CMD ["./echallan-services"]

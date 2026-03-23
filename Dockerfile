FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /veylo ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /veylo /app/veylo
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/veylo"]

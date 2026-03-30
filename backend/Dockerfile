FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download 2>/dev/null || true

COPY . .

RUN go mod tidy && CGO_ENABLED=0 go build -o serverpoint .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/serverpoint .

EXPOSE 8080

CMD ["./serverpoint"]

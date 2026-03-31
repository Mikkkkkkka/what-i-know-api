FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/api ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S app && adduser -S -G app app

COPY --from=builder /bin/api /usr/local/bin/api

ENV HTTP_ADDRESS=:8080

EXPOSE 8080

USER app

ENTRYPOINT ["/usr/local/bin/api"]

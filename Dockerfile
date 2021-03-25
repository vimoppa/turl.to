FROM golang:alpine AS builder

WORKDIR /app

RUN apk add gcc g++ ca-certificates --no-cache

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN mkdir store

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' ./cmd/turl.to

FROM scratch

WORKDIR /app

ENV PORT 8080

EXPOSE $PORT

ENTRYPOINT ["/app/turl.to"]

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/turl.to /app/turl.to
COPY --from=builder /app/store /app/store

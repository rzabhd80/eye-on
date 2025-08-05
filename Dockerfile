FROM golang:1.24-alpine AS builder


ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=off

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/app ./cmd


FROM alpine:latest

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /go/bin/app ./app
COPY --from=builder /app/migrations ./migrations

USER appuser
ENTRYPOINT ["./app"]
CMD ["api"]
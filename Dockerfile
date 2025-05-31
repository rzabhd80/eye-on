# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git


COPY go.mod go.sum ./


RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app.go


FROM alpine:latest


RUN apk --no-cache add ca-certificates


RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/


COPY --from=builder /app/main .


RUN chown appuser:appgroup main


USER appuser


EXPOSE 8080

CMD ["./main","api"]
# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./cmd/worker

# Runtime stage
FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/api /app/api
COPY --from=builder /out/worker /app/worker
COPY files/etc/transfer /app/files/etc/transfer

EXPOSE 8080

CMD ["/app/api"]

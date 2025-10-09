FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mongo-migrate .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/mongo-migrate .

# Create migrations directory (will be empty initially)
RUN mkdir -p migrations

RUN adduser -D -s /bin/sh migration
USER migration

ENTRYPOINT ["./mongo-migrate"]
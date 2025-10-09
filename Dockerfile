FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mongo-essential .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/mongo-essential .

# Create migrations directory
RUN mkdir -p migrations

COPY .env.example .
COPY examples/ examples/

RUN adduser -D -s /bin/sh migration
USER migration

ENTRYPOINT ["./mongo-essential"]

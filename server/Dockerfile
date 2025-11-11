# ========================
# Stage 1: Build Go binary
# ========================
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git bash curl tar
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/main.go

# ========================
# Stage 2: Runtime
# ========================
FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache bash curl tar
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
    -o migrate.tar.gz && \
    tar -xzf migrate.tar.gz && mv migrate /usr/local/bin/ && rm migrate.tar.gz

COPY --from=builder /app/app .
COPY migrations ./migrations
COPY entrypoint.sh .

RUN chmod +x app entrypoint.sh

EXPOSE 3001
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./app"]
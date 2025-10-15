# ---- Build Stage ----
FROM golang:1.24-alpine AS builder   

WORKDIR /app                         
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/demo/main.go

# ---- Runtime Stage ----
FROM alpine:latest                   
WORKDIR /app                          
COPY --from=builder /app/app .
COPY --from=builder /app/entrypoint.sh .

RUN chmod +x app entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./app"]

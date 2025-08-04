# ---------- Build Stage ----------

FROM --platform=linux/amd64 golang:1.24.3 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct
ENV GOMODRETRIES=10

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -o /app/app ./cmd/demo/main.go

# ---------- Runtime Stage ----------
FROM alpine:latest

# RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/entrypoint.sh .

# RUN chown appuser:appgroup /app/app /app/entrypoint.sh && \
#     chmod +x /app/app /app/entrypoint.sh

# USER appuser

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./app"]

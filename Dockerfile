# =========================
# 🧱 Builder
# =========================
FROM golang:1.25.1-alpine AS builder
ENV CGO_ENABLED=0 GO111MODULE=on
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/server

# =========================
# 🚀 Runtime kecil dg healthcheck
# =========================
FROM alpine:3.20

# Hanya paket yang benar-benar perlu
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S app && adduser -S app -G app
ENV TZ=Asia/Jakarta

WORKDIR /app
COPY --from=builder /app/server /usr/local/bin/server

USER app

EXPOSE 8081

# Healthcheck pakai curl dan port yang sesuai (8081)
HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
  CMD curl -fsS http://localhost:8081/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/server"]

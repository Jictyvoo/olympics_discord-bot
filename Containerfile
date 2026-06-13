### Stage 1: Build ###
FROM golang:1.26-bookworm AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src

# Cache module downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -trimpath \
    -o /out/olhojogo \
    ./cmd/olhojogo

### Stage 2: Runtime ###
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/olhojogo /bin/olhojogo

# Data directory for SQLite file; mount a volume here in production.
VOLUME ["/app/data"]

ENV OLH_DB_DRIVER=sqlite
ENV OLH_DB_DSN=/app/data/olhojogo.db
ENV OLH_DB_MIGRATE=false
ENV OLH_CACHE_BACKEND=file
ENV OLH_CACHE_PATH=/app/data/cache

EXPOSE 8080

ENTRYPOINT ["/bin/olhojogo"]
CMD ["serve"]

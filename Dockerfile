# Step 1: Start
FROM golang:1.22-alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app .

# Step 3: Final
FROM scratch
COPY --from=builder /app/config.docker.yaml config.yaml
COPY --from=builder /app/migrations /migrations
# COPY --from=builder /app/cert /cert
COPY --from=builder /bin/app /app
CMD ["/app"]

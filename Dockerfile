FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

COPY ./ /app/
# Build static binary with stripped debug info
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o apkingo ./cmd/apkingo/

FROM alpine:3.24

# Install CA certificates for HTTPS requests (VT/Koodous)
RUN apk --no-cache add ca-certificates \
    && addgroup -S -g 10001 apkingo \
    && adduser -S -D -H -u 10001 -G apkingo apkingo

WORKDIR /mnt/
COPY --from=builder --chown=apkingo:apkingo /app/apkingo /usr/local/bin/apkingo
USER apkingo

# Set entrypoint
ENTRYPOINT ["apkingo"]

FROM golang:alpine as builder

WORKDIR /app

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod tidy

COPY ./ /app/
# Build static binary with stripped debug info
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o apkingo ./cmd/apkingo/

FROM alpine:latest

# Install CA certificates for HTTPS requests (VT/Koodous)
RUN apk --no-cache add ca-certificates

WORKDIR /mnt/
COPY --from=builder /app/apkingo /usr/local/bin/apkingo

# Set entrypoint
ENTRYPOINT ["apkingo"]

FROM golang:alpine as builder

WORKDIR /app

COPY ./go.mod /app/go.mod
COPY ./go.mod /app/go.sum
RUN go mod download

COPY ./ /app/
RUN go build -o apkingo ./cmd/apkingo/

FROM alpine:latest
WORKDIR /mnt/
COPY --from=builder /app/apkingo ../apkingo

ENTRYPOINT ["/apkingo"]

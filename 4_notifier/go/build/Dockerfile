FROM golang:1.25-alpine AS build

# needed for cgo
RUN apk add --no-cache build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o /out/notifier ./cmd


FROM alpine:3.23

# ca-certificates: Discord and Telegram are called over HTTPS, and without the root certificates
# every send fails with "x509: certificate signed by unknown authority".
RUN apk add --no-cache ca-certificates

# Nothing here needs to be root. The user owns /data so SQLite can write its file there.
RUN adduser -D -u 10001 notifier \
    && mkdir /data \
    && chown notifier:notifier /data

COPY --from=build /out/notifier /usr/local/bin/notifier

USER notifier
WORKDIR /data

ENV NOTIFIER_PORT=8083 \
    NOTIFIER_DB_PATH=/data/notifier.db \
    NATS_URL=nats://nats:4222 \
    SMTP_HOST=mailpit \
    SMTP_PORT=1025 \
    NOTIFIER_MAIL_FROM=veille@example.org

VOLUME /data
EXPOSE 8083

ENTRYPOINT ["notifier"]

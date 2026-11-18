FROM golang:1.26-alpine AS build

RUN apk add --no-cache build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o /out/ingestor ./cmd


FROM alpine:3.23

# ca-certificates because the feeds are read over HTTPS
RUN apk add --no-cache ca-certificates

RUN adduser -D -u 10001 ingestor \
    && mkdir /data \
    && chown ingestor:ingestor /data

COPY --from=build /out/ingestor /usr/local/bin/ingestor

USER ingestor
WORKDIR /data

ENV INGESTOR_PORT=8080 \
    INGESTOR_DB_PATH=/data/ingestor.db \
    NATS_URL=nats://nats:4222

VOLUME /data
EXPOSE 8080

ENTRYPOINT ["ingestor"]

FROM golang:1.26-alpine AS build

# needed for cgo
RUN apk add --no-cache build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o /out/publisher ./cmd


FROM alpine:3.23

RUN apk add --no-cache ca-certificates

RUN adduser -D -u 10001 publisher \
    && mkdir /data \
    && chown publisher:publisher /data

COPY --from=build /out/publisher /usr/local/bin/publisher

USER publisher
WORKDIR /data

ENV PUBLISHER_PORT=8081 \
    PUBLISHER_DB_PATH=/data/publisher.db \
    NATS_URL=nats://nats:4222

VOLUME /data
EXPOSE 8081

ENTRYPOINT ["publisher"]

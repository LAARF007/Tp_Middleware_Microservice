FROM node:24-slim AS build

RUN apt-get update \
    && apt-get install -y python3 make g++ \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY package.json package-lock.json ./

# clean install
RUN npm ci --omit=dev


FROM node:24-slim

WORKDIR /app

# node:24-slim already has a "node" user
RUN mkdir /data && chown node:node /data

COPY --from=build /src/node_modules ./node_modules
COPY package.json ./
COPY src ./src

USER node

ENV NODE_ENV=production \
    NOTIFIER_PORT=8083 \
    NOTIFIER_DB_PATH=/data/notifier.db \
    NATS_URL=nats://nats:4222 \
    SMTP_HOST=mailpit \
    SMTP_PORT=1025 \
    NOTIFIER_MAIL_FROM=veille@example.org

VOLUME /data
EXPOSE 8083

CMD ["node", "src/server.js"]

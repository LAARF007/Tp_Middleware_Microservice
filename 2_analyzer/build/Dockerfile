FROM python:3.13-slim AS build

RUN python -m venv /venv
ENV PATH=/venv/bin:$PATH

WORKDIR /src

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt


FROM python:3.13-slim

RUN useradd --uid 10001 --create-home analyzer

COPY --from=build /venv /venv

WORKDIR /app
COPY analyzer ./analyzer

COPY configuration.json ./

USER analyzer

ENV PATH=/venv/bin:$PATH \
    PYTHONUNBUFFERED=1 \
    ANALYZER_CONFIG_DIR=/app \
    ANALYZER_PORT=8082 \
    NATS_URL=nats://nats:4222

# PYTHONUNBUFFERED : without it Python holds the logs in a buffer and
# `docker logs` shows nothing until the buffer fills, which looks exactly like a service that hung.

EXPOSE 8082

CMD ["python", "-m", "analyzer"]

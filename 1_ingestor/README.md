# Ingestor

## Run

You need **Go 1.26 or newer**, **gcc** (the SQLite driver `mattn/go-sqlite3` uses cgo)  
From the repository root:
```
docker compose up -d
```

Then, in this directory. Tidy / download modules :
```
go mod tidy
```

Build & run for dev :
```
go run cmd/main.go
```

## In Docker

Potentially, if not defined:
```
 docker network create --attachable veille-network
```

Then, for this service :
```
docker build -f build/Dockerfile -t veille-ingestor .
docker run --rm -p 8080:8080 -v veille_ingestor_data:/data \
  --network veille-network veille-ingestor
```
The SQLite file lives in `/data`, which has to be mounted to survive a restart.

## Configuration

Environment variables, all optional:

| Variable | Default | What it is                                                                                 |
|---|---|--------------------------------------------------------------------------------------------|
| `INGESTOR_PORT` | `8080` | HTTP port of the API                                                                       |
| `INGESTOR_DB_PATH` | `ingestor.db` | SQLite file, created if missing                                                            |

## API documentation

Documentation is visible in **api** directory ([here](api/swagger.json)). To regenerate it after
changing an annotation:
```
go run github.com/swaggo/swag/cmd/swag@latest init -g ./internal/swagger_expose.go -o ./api
```

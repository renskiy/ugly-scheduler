# ugly-scheduler
Ugly-scheduler is just a simple schedule service. It uses [Twirp-based](https://github.com/twitchtv/twirp) API which provides benefits from Google's [protobuf](https://en.wikipedia.org/wiki/Protocol_Buffers) (code generation) and REST JSON API (simple request construction).

## Requirements
* Go 1.13
* PostgreSQL 9.5 or greater

## Configuration
Service configured using command-line options and environment variables.

### Command-line options
* `--addr`, addr:port to run service on (default: `:8001`)

### Environment variables
* `DB_HOST`, DB host (default: `localhost`)
* `DB_PORT`, DB port (default: `5432`)
* `DB_NAME`, DB name (default: `ugly-scheduler`)
* `DB_USER`, DB user (default: `postgres`)
* `DB_PASSWORD`, DB password (default not set)

## Environment setup
1. clone this repository to `$GOPATH/src/github.com/renskiy` directory
2. install or update `protoc-gen-go` package

### Requirements
    go mod download
    
## Database
    createdb ugly-scheduler

### Applying migrations
    go get -u github.com/rubenv/sql-migrate/...
    go run ./cmd/sql-migrate-config | sql-migrate up -config /dev/stdin -env default

## Tests
Create special DB for tests, set `DB_NAME_TEST`* env var and run following command:
    
    go test -v -count 1 -p 1 ./...

\* *default value for `DB_NAME_TEST` is `ugly-scheduler-test`*

## Run
    go run ./cmd/server --addr :8001
    go run ./cmd/server --addr :8002

### Making requests
With default settings one can do requests like following:

    curl -X POST "http://localhost:8001/twirp/scheduler.Scheduler/Schedule" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"delay\": \"1\", \"message\": \"Hello world\"}"

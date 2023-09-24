# DashboardX-API-PoC

## Setup

1. `docker compose run app go run ./cmd/cli/main.go create`
2. `docker compose run app go run ./cmd/cli/main.go up`
3. `docker compose run app go run ./cmd/cli/main.go seed`
4. `docker compose up -d`

## Testing

In docker:
- `docker compose run app go test -v ./...`

Locally:
- `go test -v ./...`

## Swagger

Available at the link: `<HOST>/swagger/index.html`

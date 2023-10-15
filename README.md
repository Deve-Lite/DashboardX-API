# DashboardX-API

## Setup

Run `./scripts/d_setup.ps1`, it setups whole system from the scratch. Be aware that all stored data will be removed.

## Upgrade

To upgrade (in case of new changes or migrations) the API run `./scripts/d_upgrade.ps1` 

## Testing

In docker:
- `docker compose run app go test -v ./...`

Locally:
- `go test -v ./...`

## Swagger

Available at the link: `<HOST>/swagger/index.html`

# Upgrades version of the API and runs all pending migrations

docker compose stop app

docker compose run --rm app go run ./cmd/cli/main.go -op=up

docker compose rm -vfs app
docker rmi -f dashboardx-api-app

docker compose build --no-cache app

docker compose up -d

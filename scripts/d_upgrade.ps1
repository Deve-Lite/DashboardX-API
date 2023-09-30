# Upgrades version of the API and runs all pending migrations

docker compose stop app

docker compose run --rm app go run ./cmd/cli/main.go up

docker compose rm -vfs app
docker rmi -f dashboardx-app-app

docker compose build --no-cache app

docker compose up -d

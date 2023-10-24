# Clean setup of the system, forces to delete all stored data!

param (
    [Parameter(HelpMessage="Force to prune docker system")]
    [switch]$Prune = $false
)

docker compose rm -vfs
docker rmi -f dashboardx-api-app

docker volume rm dashboardx-api_postgres-data
docker volume rm dashboardx-api_redis-data
docker volume rm dashboardx-api_smtp4dev-data

if ($Prune) {
    docker system prune --all --force
} 

docker compose build --no-cache

docker compose run --rm app go run ./cmd/cli/main.go create
docker compose run --rm app go run ./cmd/cli/main.go up
docker compose run --rm app go run ./cmd/cli/main.go seed

docker compose up -d

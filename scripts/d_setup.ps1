# Clean setup of the system, forces to delete all stored data!

docker compose rm -vfs

docker compose run --rm app go run ./cmd/cli/main.go create
docker compose run --rm app go run ./cmd/cli/main.go up
docker compose run --rm app go run ./cmd/cli/main.go seed

docker compose up -d

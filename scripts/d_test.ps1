# Runs integration tests in the Docker

docker compose run --rm app go test -v -timeout 30s ./internal/interfaces/http/rest/handler

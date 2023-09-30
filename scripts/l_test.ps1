# Runs integration tests locally

docker compose up -d # Postgres & Redis need to be up for tests

go test -v -timeout 30s ./internal/interfaces/http/rest/handler

# Runs the API locally
dev:
	@go run ./cmd/server/main.go


# Runs integration tests locally
# Postgres & Redis need to be up for tests
dev-test:
	@docker compose up -d
	
	@go test -v -timeout 30s \
		./internal/interfaces/http/rest/handler \
		./internal/application


# Clean setup of the system, forces to delete all stored data!
setup:
	@docker compose rm -vfs
	@docker rmi -f dashboardx-api-app

	@docker compose build --no-cache

	@docker compose run --rm app go run ./cmd/cli/main.go create
	@docker compose run --rm app go run ./cmd/cli/main.go up
	@docker compose run --rm app go run ./cmd/cli/main.go seed

	@docker compose up -d


# The setup with extra prune
setup-prune:
	@docker system prune --all --force
	@$(MAKE) setup


# Runs integration tests in the Docker
test:
	@docker compose run --rm app go test -v -timeout 30s \
    	./internal/interfaces/http/rest/handler \
    	./internal/application


# Upgrades version of the API and runs all pending migrations
upgrade:
	@docker compose stop app

	@docker compose run --rm app go run ./cmd/cli/main.go up

	@docker compose rm -vfs app
	@docker rmi -f dashboardx-api-app

	@docker compose build --no-cache app

	@docker compose up -d


# Generates Swagger documention
# Replace 'x-nullable' with 'nullable' to proper display it in the Swagger
swagger:
	@swag fmt
	@swag init -g ./cmd/server/main.go

	@sed -i 's/x-nullable/nullable/g' ./docs/docs.go
	@sed -i 's/x-nullable/nullable/g' ./docs/swagger.json
	@sed -i 's/x-nullable/nullable/g' ./docs/swagger.yaml
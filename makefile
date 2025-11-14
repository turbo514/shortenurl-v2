prod-infra-up:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-infra.yaml up
prod-infra-down:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-infra.yaml down
prod-app-up:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-app.yaml up
prod-app-down:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-app.yaml down -v
dev:
	cd deploy/docker_compose && docker compose -p my-dev-compose -f compose-dev.yaml up

buf-generate:
	buf generate

link-sqlc:
	cd internal/link && sqlc generate -f infrastructure/mysqldb/sqlc.yaml
link-run:
	cd internal/link && go run main.go

tenant-run:
	cd internal/tenant && go run main.go

analytics-run:
	cd internal/analytics && go run main.go

gateway-run:
	cd internal/gateway && go run main.go
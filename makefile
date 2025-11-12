up:
	cd deploy/docker_compose && docker compose up
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
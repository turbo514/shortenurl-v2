prod-infra-up:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-infra.yaml up -d
prod-infra-down:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-infra.yaml down
prod-app-up:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-app.yaml --env-file ./prod.env up -d
prod-app-down:
	cd deploy/docker_compose && docker compose -p my-prod-compose -f compose-prod-app.yaml down -v
dev-up:
	cd deploy/docker_compose && docker compose -p my-dev-compose -f compose-dev.yaml up
kafka-init:
	docker exec -it my-prod-compose-kafka-prod-1 /bin/bash /opt/kafka/create-topics.sh

#migrate -path ./internal/migrations/tenant_mysql/ -database "mysql://user:1234@tcp(127.0.0.1:3306)/tenant_dev" up
#migrate -path ./internal/migrations/link_mysql/ -database "mysql://user:1234@tcp(127.0.0.1:3306)/link_dev" u
# migrate -path ./internal/migrations/clickhouse/ -database "clickhouse://127.0.0.1:9000?username=user&password=1234&database=default&x-multi-statement=true" up

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
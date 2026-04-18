MYSQLDEF ?= mysqldef
MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 3306
MYSQL_USER ?= root
MYSQL_PASSWORD ?=
MYSQL_DATABASE ?= techcv_app

.PHONY: mysql-up mysql-down schema-dry-run schema-apply sqlc-generate oapi-generate

mysql-up:
	@docker compose -f ../docker-compose.yml up -d mysql

mysql-down:
	@docker compose -f ../docker-compose.yml stop mysql

sqlc-generate:
	@go tool sqlc generate

oapi-generate:
	@go tool oapi-codegen --config openapi/oapi-codegen.yaml openapi/openapi.yaml

schema-dry-run:
	@$(MYSQLDEF) \
		-u $(MYSQL_USER) \
		-p$(MYSQL_PASSWORD) \
		-h $(MYSQL_HOST) \
		-P $(MYSQL_PORT) \
		$(MYSQL_DATABASE) \
		--dry-run < migrations/schema.sql

schema-apply:
	@$(MYSQLDEF) \
		-u $(MYSQL_USER) \
		-p$(MYSQL_PASSWORD) \
		-h $(MYSQL_HOST) \
		-P $(MYSQL_PORT) \
		$(MYSQL_DATABASE) < migrations/schema.sql

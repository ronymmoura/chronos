include .env

MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DB_URL := $(CONNECTION_STRING)

migration-new:
	docker run --rm -v $(MAKEFILE_DIR):/src -w /src --network host migrate/migrate create -ext sql -dir=internal/db/migrations -seq $(name)

migrate-up:
	docker run --rm -v $(MAKEFILE_DIR):/src -w /src --network host migrate/migrate -path=internal/db/migrations/ -database="$(DB_URL)" up

.PHONY:
	migration-new migrate-up

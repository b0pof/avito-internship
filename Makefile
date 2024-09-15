
ifneq ("$(wildcard .env)","")
include .env
endif

.PHONY: setup
setup: ## Установить все необходимые утилиты
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrations-up
migrations-up: ## Накатить миграции
	goose -dir migrations postgres $(POSTGRES_CONN) up

.PHONY: migrations-down
migrations-down: ## Откатить миграции
	goose -dir migrations postgres $(POSTGRES_CONN) down

.PHONY: migration-create
migration-create: ## Пример команды для создания миграции
	@echo "goose -dir migrations create <add_some_column> sql"

.PHONY: run-app
run-app: ## Запустить приложение
	go run cmd/main.go

.PHONY: run-db
run-db: ## Поднять базу
	docker compose -f docker-compose.yaml up -d

.PHONY: stop-db
stop-all: ## Остановить базу
	docker compose -f docker-compose.yaml down

.PHONY: build
build: ## Сбилдить исполняемый файл приложения
	go build -o ./bin/app ./cmd/main.go

.PHONY: lint
lint: ## Проверить код линтером
	golangci-lint run ./... -c golangci.yml

.PHONY: clean
clean: ## Удалить временные файлы
	rm -f ./bin/app

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:=help

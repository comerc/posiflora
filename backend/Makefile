.PHONY: help test mockery run build test-order

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Запустить тесты
	go test -v ./...

test-order: ## Создать тестовый заказ через API
	@./scripts/test-order.sh

mockery: ## Генерировать моки
	@echo "Генерация моков через mockery..."
	@if command -v mockery > /dev/null; then \
		mockery --all --output=./internal/mocks; \
	else \
		echo "mockery не установлен. Установите: go install github.com/vektra/mockery/v2@latest"; \
	fi

run: ## Запустить сервер
	go run cmd/server/main.go

build: ## Собрать бинарник
	go build -o bin/server cmd/server/main.go

clean: ## Очистить артефакты сборки
	rm -rf bin/


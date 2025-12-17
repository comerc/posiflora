.PHONY: help test test-order run build clean install setup db-migrate db-reset

help: ## Показать справку
	@echo "Доступные команды для backend2 (PHP Symfony):"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Запустить тесты
	php bin/phpunit

test-order: ## Создать тестовый заказ через API
	@if [ -f "./scripts/test-order.sh" ]; then \
		echo "Создание тестового заказа..."; \
		API_URL=http://localhost:3000 SHOP_ID=$(SHOP_ID) NUMBER=$(NUMBER) TOTAL=$(TOTAL) CUSTOMER=$(CUSTOMER) ./scripts/test-order.sh; \
	else \
		echo "Скрипт test-order.sh не найден в ./scripts/"; \
		echo "Убедитесь, что backend запущен (make run-symfony)"; \
	fi

run: ## Запустить development сервер
	php -S localhost:3000 -t public

run-symfony: ## Запустить сервер через Symfony CLI
	symfony server:start --port=3000 --no-tls

build: ## Очистить кэш и подготовить для production
	composer install --no-dev --optimize-autoloader
	php bin/console cache:clear --env=prod
	php bin/console cache:warmup --env=prod

clean: ## Очистить кэш и временные файлы
	rm -rf var/cache/*
	rm -rf var/log/*
	rm -rf var/sessions/*

install: ## Установить зависимости
	composer install

setup: install db-migrate ## Полная настройка проекта

db-migrate: ## Выполнить миграции базы данных
	php bin/console doctrine:migrations:migrate --no-interaction

db-reset: ## Сбросить базу данных и выполнить миграции заново
	@if [ -f "var/data/data_dev.db" ]; then rm var/data/data_dev.db; fi
	@if [ -f "var/data/data_prod.db" ]; then rm var/data/data_prod.db; fi
	php bin/console doctrine:migrations:migrate --no-interaction

cache-clear: ## Очистить кэш
	php bin/console cache:clear

cache-warmup: ## Прогреть кэш
	php bin/console cache:warmup

docker-up: ## Запустить через Docker Compose
	docker-compose up -d

docker-down: ## Остановить Docker Compose
	docker-compose down

docker-build: ## Собрать Docker образы
	docker-compose build --no-cache

docker-logs: ## Посмотреть логи Docker
	docker-compose logs -f

# Posiflora Backend v2 (PHP Symfony)

Вторая версия backend для интеграции Telegram с магазинами Posiflora. Реализована на PHP с использованием Symfony framework.

## Архитектура

Проект следует стандартной Symfony архитектуре:

- **Entities** (`src/Entity/`) - модели данных Doctrine ORM
- **Services** (`src/Service/`) - бизнес-логика
- **Controllers** (`src/Controller/`) - HTTP контроллеры API
- **Repositories** (`src/Repository/`) - доступ к данным
- **Telegram** (`src/Telegram/`) - клиенты для работы с Telegram API

## API Endpoints

### Подключение Telegram-интеграции

```
POST /shops/{shopId}/telegram/connect
```

```json
{
  "botToken": "123456:ABC-DEF...",
  "chatId": "987654321",
  "enabled": true
}
```

### Создание заказа

```
POST /shops/{shopId}/orders
```

```json
{
  "number": "A-1005",
  "total": 2490.0,
  "customerName": "Анна"
}
```

### Статус интеграции

```
GET /shops/{shopId}/telegram/status
```

```json
{
  "enabled": true,
  "chat_id": "***54321",
  "last_sent_at": "2025-01-01T12:00:00Z",
  "sent_count": 5,
  "failed_count": 1
}
```

## Запуск

### Локальная разработка

1. **Установка и настройка:**

   ```bash
   make setup  # Установка зависимостей + миграции БД
   ```

2. **Запуск сервера:**
   ```bash
   make run    # Development сервер на localhost:3000
   ```

Или через Symfony CLI:

```bash
make run-symfony  # Symfony CLI на порту 3000
```

Сервер будет доступен на `http://localhost:3000`

### Docker

1. **Сборка и запуск:**
   ```bash
   docker-compose up --build
   ```

Сервер будет доступен на `http://localhost:8081`

### Переменные окружения

- `APP_ENV` - окружение (dev/prod)
- `DATABASE_URL` - URL базы данных
- `TELEGRAM_MOCK_MODE` - использовать мок Telegram клиента (true/false)
- `CORS_ALLOW_ORIGIN` - разрешенные origins для CORS

## Тестирование

### Запуск тестов

```bash
php bin/phpunit
```

### Запуск конкретного теста

```bash
php bin/phpunit tests/Service/OrderServiceTest.php
```

## Сиды данных

Для загрузки тестовых данных:

```bash
php bin/console doctrine:migrations:migrate
```

Тестовые данные включают:

- Магазин "Тестовый магазин Posiflora" (ID: 1)
- 8 тестовых заказов

## Telegram интеграция

### Mock режим (по умолчанию)

При `TELEGRAM_MOCK_MODE=true` сообщения логируются в консоль вместо отправки в Telegram.

### Реальный режим

При `TELEGRAM_MOCK_MODE=false` сообщения отправляются в реальный Telegram API.

Для получения `chat_id`:

1. Добавьте бота в канал/группу
2. Отправьте сообщение в канал
3. Вызовите `https://api.telegram.org/bot{token}/getUpdates`
4. Найдите `chat.id` в ответе

## Структура БД

- `shops` - магазины
- `telegram_integrations` - настройки Telegram интеграции
- `orders` - заказы
- `telegram_send_log` - лог отправки уведомлений

## Особенности реализации

- **Идемпотентность**: повторные вызовы создания заказа не создают дубли в логе отправки
- **Отказоустойчивость**: ошибки Telegram API не препятствуют созданию заказов
- **Маскирование**: chat_id маскируется в ответах API (параметр `mask=disabled` отключает маскирование)
- **CORS**: настроена поддержка кросс-доменных запросов

## Разработка

### Создание новых entities

```bash
php bin/console make:entity EntityName
```

### Создание миграций

```bash
php bin/console doctrine:migrations:diff
php bin/console doctrine:migrations:migrate
```

### Создание контроллеров

```bash
php bin/console make:controller ControllerName
```

### Makefile команды

Проект включает Makefile для удобного управления:

```bash
make help          # Показать все доступные команды
make test          # Запустить тесты
make run           # Запустить development сервер (localhost:3000)
make run-symfony   # Запустить через Symfony CLI
make db-migrate    # Выполнить миграции БД
make db-reset      # Сбросить и пересоздать БД
make cache-clear   # Очистить кэш
make docker-up     # Запустить через Docker
make docker-down   # Остановить Docker
make setup         # Полная настройка проекта
```

Создайте тестовый заказ через API:

```bash
# Через Makefile с параметрами
make test-order SHOP_ID=1 NUMBER=TEST-001 TOTAL=1000 CUSTOMER=Test

# Или напрямую через curl
curl -X POST http://localhost:3000/shops/1/orders \
  -H "Content-Type: application/json" \
  -d '{"number":"TEST-001","total":1000,"customerName":"Test"}'
```

**Доступные параметры для make test-order:**

- `SHOP_ID` - ID магазина (по умолчанию: 1)
- `NUMBER` - номер заказа (по умолчанию: TEST-001)
- `TOTAL` - сумма заказа (по умолчанию: 1000)
- `CUSTOMER` - имя клиента (по умолчанию: Test)

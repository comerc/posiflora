# posiflora

## Инструкции по запуску

### 1. Backend

#### Через Docker

```bash
docker-compose up -d
```

Backend будет доступен на `http://localhost:8080`

#### Без Docker

1. Установите PostgreSQL и создайте БД:

```bash
createdb posiflora
```

2. Создайте `.env` файл в `backend/` (скопируйте `backend/.env.example` в `backend/.env` и при необходимости измените значения)

3. Запустите сервер:

```bash
cd backend
make run
```

### 2. Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend будет доступен на `http://localhost:5173`

### 3. Тестовые данные

Тестовые данные создаются автоматически при запуске сервера через Docker (миграции применяются автоматически):

- 1 магазин (ID=1, название "Тестовый магазин Posiflora")
- 8 тестовых заказов (A-1001 до A-1008)

Либо установите [goose](https://github.com/pressly/goose) и выполните миграцию данных, если backend без Docker.

> Если магазин отсутствует при сохранении интеграции, то будет добавлен автоматически.

### 4. Тесты

```bash
cd backend
make test
```

### 5. Telegram режим

По умолчанию используется **мок-режим** (сообщения логируются в консоль).

**Как проверить мок-режим:**

Зайдите на станицу интеграции: `http://localhost:5173/shops/1/growth/telegram`

1. Подключите интеграцию через UI (введите любые тестовые значения для botToken и chatId)
2. Включите интеграцию (переключатель "Включено"), нажмите "Сохранить"
3. Создайте тестовый заказ через API:
   ```bash
   curl -X POST http://localhost:8080/shops/1/orders \
     -H "Content-Type: application/json" \
     -d '{"number":"TEST-001","total":1000,"customerName":"Test"}'
   ```
   или `make test-order SHOP_ID=1 NUMBER=TEST-001 TOTAL=1000 CUSTOMER=Test`
4. Проверьте логи backend - должно появиться сообщение `[MOCK TELEGRAM] Sending message...`
5. В реальном режиме сообщение будет отправлено через Telegram Bot API
6. Перезагрузите страницу интеграции, чтобы увидеть обновлённый статус интеграции (последняя отправка).

Для включения реальной отправки установите в `backend/.env`:

```env
TELEGRAM_MOCK_MODE=false
```

### 6. Допущения и упрощения

1. **Telegram в мок-режиме по умолчанию** - для упрощения разработки и тестирования
2. **PostgreSQL вместо SQLite** - для соответствия production-окружению
3. **Минимальная валидация** - только обязательные поля
4. **Нет аутентификации** - для MVP не требуется
5. **Простой error handling** - базовые проверки ошибок

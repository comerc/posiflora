# План выполнения технического задания

## Этап 1: Настройка проекта и инфраструктуры

- [x] Выбрать стек технологий (Backend: PHP/Symfony или Go, Frontend: React + TypeScript)
- [x] Инициализировать backend проект
- [x] Инициализировать frontend проект (React + TypeScript)
- [x] Настроить структуру проекта (разделение на слои)
- [x] Создать Dockerfile для backend
- [x] Создать docker-compose.yml (backend + database)
- [x] Настроить переменные окружения (.env файлы)
- [x] Настроить базовую конфигурацию проекта

## Этап 2: База данных

- [x] Создать схему БД:
  - [x] Таблица `shops` (id, name)
  - [x] Таблица `telegram_integrations` (id, shop_id, bot_token, chat_id, enabled, created_at, updated_at, unique(shop_id))
  - [x] Таблица `orders` (id, shop_id, number, total, customer_name, created_at)
  - [x] Таблица `telegram_send_log` (id, shop_id, order_id, message, status, error, sent_at, unique(shop_id, order_id))
- [x] Создать миграции/схему БД
- [x] Настроить подключение к БД
- [x] Создать seed данные:
  - [x] Минимум 1 магазин
  - [x] 5-10 тестовых заказов

## Этап 3: Backend - Архитектура и слои

- [x] Создать интерфейс TelegramClient
- [x] Реализовать TelegramClient (реальный или мок)
- [x] Создать репозитории для работы с БД:
  - [x] ShopRepository
  - [x] TelegramIntegrationRepository
  - [x] OrderRepository
  - [x] TelegramSendLogRepository
- [x] Создать сервисный слой:
  - [x] TelegramIntegrationService
  - [x] OrderService
  - [x] TelegramNotificationService

## Этап 4: Backend API - Эндпоинты

### 4.1. Подключение Telegram-интеграции

- [x] Реализовать POST `/shops/{shopId}/telegram/connect`
- [x] Валидация payload (botToken, chatId не пустые)
- [x] Реализовать upsert логику (обновление существующей интеграции)
- [x] Обновление updated_at
- [x] Возврат сохранённой конфигурации

### 4.2. Создание заказа

- [x] Реализовать POST `/shops/{shopId}/orders`
- [x] Создание заказа в БД
- [x] Проверка наличия активной Telegram-интеграции
- [x] Отправка уведомления в Telegram (если интеграция включена)
- [x] Реализация идемпотентности (проверка telegram_send_log)
- [x] Логирование попытки отправки (SENT/FAILED)
- [x] Обработка ошибок (заказ создаётся даже при ошибке Telegram)
- [x] Возврат заказа + статус отправки (sent/failed/skipped)

### 4.3. Статус интеграции

- [x] Реализовать GET `/shops/{shopId}/telegram/status`
- [x] Возврат enabled
- [x] Возврат chatId (с частичным маскированием)
- [x] Возврат lastSentAt
- [x] Подсчёт sentCount за 7 дней
- [x] Подсчёт failedCount за 7 дней

## Этап 5: Telegram интеграция

- [x] Реализовать TelegramClient с интерфейсом
- [x] Интеграция с Telegram Bot API:
  - [x] POST запрос к `https://api.telegram.org/bot{token}/sendMessage`
  - [x] Обработка ответа API
  - [x] Обработка ошибок
- [x] Реализовать мок-режим для разработки
- [x] Настроить переключение между реальным и мок-режимом
- [x] Форматирование сообщения: "Новый заказ {number} на сумму {total} ₽, клиент {customerName}"

## Этап 6: Frontend

- [x] Vite+React+TS
- [x] import via alias `@`
- [x] vite-env.d.ts replaced to `"types": ["vite/client"],`
- [x] .prettierrc.json
- [x] prettier-plugin-tailwindcss
- [x] postcss config: autoprefixer + cssnano
- [x] Tailwind
- [x] Antd
- [x] customisation of Ant Design with Tailwind
- [x] fix https://ant.design/docs/react/v5-for-19
- [x] Настроить роутинг (React Router)
  - [x] layouts
  - [x] NotFoundPage
- [x] Создать страницу `/shops/:shopId/growth/telegram`
- [x] Реализовать форму подключения:
  - [x] Поле ввода botToken
  - [x] Поле ввода chatId
  - [x] Тумблер enabled
  - [x] Кнопка "Сохранить" (POST /telegram/connect)
- [x] Реализовать блок статуса:
  - [x] Отображение enabled
  - [x] Отображение lastSentAt
  - [x] Отображение sent/failed за 7 дней
  - [x] Автообновление статуса (GET /telegram/status)
- [x] Добавить подсказку "как узнать chatId"
- [x] Обработка ошибок и валидация форм
- [x] Базовые стили (минимальный UI)

## Этап 7: Тесты

- [x] Настроить тестовое окружение
- [x] Тест 1: При создании заказа и включённой интеграции вызывается TelegramClient и пишется лог SENT
- [x] Тест 2: Повторная отправка/повторное создание заказа не создаёт дублей telegram_send_log и не шлёт сообщение повторно (идемпотентность)
- [x] Тест 3: При ошибке TelegramClient пишется лог FAILED, а заказ всё равно создаётся
- [x] Дополнительные тесты (опционально):
  - [x] Тест валидации подключения интеграции
  - [x] Тест статуса интеграции

## Этап 8: Документация и финализация

- [x] Обновить README.md:
  - [x] Инструкция по запуску backend (docker/без docker)
  - [x] Инструкция по запуску frontend
  - [x] Инструкция по сидению тестовых данных
  - [x] Инструкция по запуску тестов
  - [x] Описание реальной Telegram-отправки или мок-режима
  - [x] Список допущений/упрощений
- [x] Проверить все требования из ТЗ
- [x] Проверить идемпотентность
- [x] Проверить обработку ошибок
- [x] Проверить разделение слоёв
- [x] Финальное тестирование всего функционала

## Этап 9: Подготовка к отправке

- [x] Создать репозиторий (GitHub/GitLab) или подготовить архив
- [x] Убедиться, что все файлы включены
- [x] Проверить, что проект запускается из коробки
- [x] Подготовить краткое описание допущений/упрощений

## Дополнительные заметки

### Важные моменты:

- Идемпотентность: использовать unique(shop_id, order_id) в telegram_send_log
- Ошибки Telegram не должны ронять создание заказа
- Разделение слоёв: сервис/хендлер → репозиторий → TelegramClient
- Безопасность: токены в env, маскирование в статусе

### Возможные упрощения:

- Можно использовать SQLite вместо PostgreSQL/MySQL для упрощения
- Минимальный UI без сложного дизайна
- Мок-режим для Telegram (с возможностью переключения на реальный)

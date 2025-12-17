FROM php:8.4-fpm-alpine

# Устанавливаем зависимости
RUN apk add --no-cache \
    git \
    unzip \
    libzip-dev \
    sqlite-dev \
    && docker-php-ext-install pdo_sqlite zip

# Устанавливаем Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Копируем код
WORKDIR /var/www/html
COPY . .

# Устанавливаем зависимости
RUN composer install --no-dev --optimize-autoloader

# Создаем директории для кэша и логов
RUN mkdir -p var/cache var/log var/data \
    && chown -R www-data:www-data var/

# Переключаемся на пользователя www-data
USER www-data

EXPOSE 9000

CMD ["php-fpm"]

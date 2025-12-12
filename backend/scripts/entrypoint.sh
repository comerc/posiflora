#!/bin/sh
set -e

echo "Waiting for database..."
until pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}"; do
  sleep 1
done

echo "Running migrations..."
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}"
goose -dir /app/migrations up

echo "Starting server..."
exec ./server


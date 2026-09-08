# Music Service — микросервисы на Go

Минимальный музыкальный сервис по ER-схеме: пользователи, каталог (альбомы, исполнители, треки), подписки.

## Архитектура

| Сервис | Порт | Назначение |
|--------|------|------------|
| **gateway** | 8080 | HTTP API, маршрутизация через gRPC |
| **user-service** | 50051 | Регистрация, логин (JWT), CRUD пользователей |
| **catalog-service** | 50052 | CRUD альбомов, исполнителей, треков |
| **subscription-service** | 50053 | CRUD подписок на треки |
| **postgres** | 5432 | База данных |
| **rabbitmq** | 5672 / 15672 | Брокер событий |

### Асинхронные события

При удалении пользователя (`DELETE /api/users/{id}`) user-service публикует событие `user.deleted` в RabbitMQ:

- **catalog-service** — скрывает треки из библиотеки пользователя (`hidden = true`)
- **subscription-service** — удаляет все подписки пользователя

## Запуск

```bash
docker compose up --build
```

Проверка:

```bash
curl http://localhost:8080/health
```

## Примеры API

### Регистрация

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@mail.com","password":"secret","phone_number":"79001234567","country":"RU"}'
```

### Логин

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@mail.com","password":"secret"}'
```

Сохраните `token` из ответа и передавайте в заголовке `Authorization: Bearer <token>`.

### Создать альбом

```bash
curl -X POST http://localhost:8080/api/albums \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"album_name":"Abbey Road","release_year":"1969","genre":"Rock"}'
```

### Создать подписку

```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"users_id":"00001","track_list_id":"00001","start_date":"2026-01-01","end_date":"2026-12-31"}'
```

### Удалить пользователя (асинхронная очистка данных)

```bash
curl -X DELETE http://localhost:8080/api/users/00001 \
  -H "Authorization: Bearer <token>"
```

## Структура проекта

```
cmd/
  gateway/              # HTTP → gRPC
  user-service/         # JWT, users, my_music
  catalog-service/      # albums, performers, track_list
  subscription-service/ # subscription_validaty
proto/                  # gRPC контракты
gen/                    # gRPC stubs (JSON codec)
internal/               # events, jwt, db, grpc helpers
deploy/init.sql         # схема PostgreSQL
docker-compose.yml
```

## Локальная разработка

```bash
go mod tidy
go run ./cmd/user-service
go run ./cmd/catalog-service
go run ./cmd/subscription-service
go run ./cmd/gateway
```

Требуются запущенные PostgreSQL и RabbitMQ (можно поднять только их через `docker compose up postgres rabbitmq`).

# Wishlist Service

REST API сервис вишлистов с чистой архитектурой.

Функционал:
- регистрация и логин по email + password (JWT);
- CRUD вишлистов для авторизованного пользователя;
- CRUD подарков внутри вишлиста;
- получение списка подарков в конкретном вишлисте;
- публичный доступ к вишлисту по токену;
- публичное бронирование подарка (без авторизации).
- автоматический запуск миграций через `golang-migrate` при старте приложения.

## Архитектура

```text
docs/
  openapi/         // OpenAPI спецификация и handler раздачи
internal/
  domain/          // бизнес-сущности
  errs/            // общие ошибки приложения
  application/     // use case
  adapters/http/   // роуты, middleware, handlers
  adapters/repository/postgres/ // репозитории PostgreSQL
  bootstrap/       // сборка зависимостей
  infrastructure/  // config, auth, db
```

## Запуск

```bash
docker-compose up --build
```

Сервис доступен на `http://localhost:8080`.

OpenAPI спецификация доступна на `http://localhost:8080/openapi.yaml`.

## Примеры запросов

### 1. Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 2. Логин

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 3. Создание вишлиста

```bash
curl -X POST http://localhost:8080/api/v1/wishlists/ \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "event_title":"День рождения",
    "description":"Подарки на ДР",
    "event_date":"2026-08-15T00:00:00Z"
  }'
```

### 4. Публичный просмотр по токену

```bash
curl http://localhost:8080/api/v1/public/<PUBLIC_TOKEN>
```

### 5. Список подарков в вишлисте

```bash
curl http://localhost:8080/api/v1/wishlists/<WISHLIST_ID>/items/ \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 6. Публичное бронирование подарка

```bash
curl -X POST http://localhost:8080/api/v1/public/<PUBLIC_TOKEN>/reserve/<ITEM_ID>
```

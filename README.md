# GO_Auth

Легковесный Go-проект с авторизацией, JWT-токенами, ролевой моделью и управлением правами доступа.

## Обзор

- Веб-сервер: `gin`.
- База данных: `sqlite` (`gorm` ORM).
- JWT для аутентификации (`golang-jwt/jwt/v5`).
- Хэширование паролей: `bcrypt`.
- Модель ролей/разрешений: `admin`, `user`.

## Структура проекта

- `main.go` — инициализация сервера, маршрутов и middleware.
- `auth.go` — логика JWT, хэширование паролей, права доступа, инициализация БД.
- `handlers.go` — HTTP-обработчики (`register`, `login`, `profile`, `products`, `create-permission`).
- `models.go` — модели данных: `User`, `Resource`, `Permission`, запросы/ответы.

## Модель данных

### User

- `ID`, `Email`, `PasswordHash`, `FirstName`, `LastName`, `Role`, `IsActive`.
- Роль по умолчанию: `user`.

### Resource

- `Name` примеры: `products`, `orders`, `users`.

### Permission

- `Role`, `ResourceID`, флаги: `CanView`, `CanCreate`, `CanEdit`, `CanDelete`.

## Запуск проекта

1. Убедитесь, что установлен [Go 1.20+]
2. В директории проекта:

```bash
cd d:/GO_Auth
go mod tidy
go run .
```

3. Сервер по умолчанию слушает `:8080`.

## Инициализация базы данных

В `InitDB()`:
- создается `test.db`.
- выполняется `AutoMigrate` для моделей.
- добавляются ресурсы: `products`, `orders`, `users`.
- создается админ с логином `admin@test.com`, паролем `admin123`.

## Маршруты API

### Публичные

- `POST /register`
  - Тело JSON: `{ "email", "password", "first_name", "last_name" }`.
  - Ответ: `201` или ошибка.

- `POST /login`
  - Тело JSON: `{ "email", "password" }`.
  - Ответ: `200` + `token`, `email`, `role`.

- `GET /`
  - `200` "Auth system on Go".

### Защищенные (Авторизация JWT)

Все маршруты требуют заголовок `Authorization: Bearer <token>`.

- `GET /profile` — информацию о пользовательском аккаунте.
- `GET /api/products` — список продуктов (проверка `permission view` по ресурсам).
- `POST /admin/create-permission` — админ устанавливает права.
    - Тело JSON: `{ "role", "resource", "can_view", "can_create", "can_edit", "can_delete" }`.

## Безопасность и поведение

- Пароли хранятся хэшами via `bcrypt`.
- Токены живут 30 минут.
- `admin` имеет полный доступ в `CheckPermission`.
- Отключенные пользователи (`IsActive=false`) получают `403`.

## Примеры запросов

### Регистрация

```bash
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"password123","first_name":"Ivan","last_name":"Ivanov"}'
```

### Вход

```bash
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"password123"}'
```

### Профиль

```bash
curl http://localhost:8080/profile -H "Authorization: Bearer <token>"
```

### Создание прав (admin)

```bash
curl -X POST http://localhost:8080/admin/create-permission -H "Authorization: Bearer <admin-token>" -H "Content-Type: application/json" -d '{"role":"user","resource":"products","can_view":true}'
```

## Улучшения (future work)

- хранение `jwtSecret` в переменных окружения;
- refresh токены;
- RBAC более высокой гибкости (сущность `Role`);
- отдельно хранение `Resource` и список объектов;
- тесты (unit + integration); 
- сборка Docker.

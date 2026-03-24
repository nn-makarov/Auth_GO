# GO_Auth

Go-сервис аутентификации и авторизации с JWT, ролевой моделью и Docker.

## Стек

- Go 1.26, Gin, GORM, SQLite
- JWT, bcrypt
- Docker

## Запуск

### Локально

```bash
go mod tidy
go run .
Docker
bash
docker run -p 8080:8080 -v $(pwd):/app -w /app golang:1.26 go run .
Сервер будет доступен на http://localhost:8080

API
Метод	Эндпоинт	Описание
POST	/register	Регистрация
POST	/login	Вход, получение JWT
GET	/profile	Информация о пользователе
GET	/api/products	Список продуктов (проверка прав)
POST	/admin/create-permission	Создание прав (только admin)
Примеры
Регистрация
bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"1234567890","first_name":"Ivan","last_name":"Ivanov"}'
Логин
bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"1234567890"}'
Профиль
bash
curl http://localhost:8080/profile \
  -H "Authorization: Bearer <token>"
Создание прав (admin)
bash
curl -X POST http://localhost:8080/admin/create-permission \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"role":"user","resource":"products","can_view":true}'
Продукты
bash
curl http://localhost:8080/api/products \
  -H "Authorization: Bearer <user_token>"
Данные по умолчанию
При первом запуске создаётся:

Администратор: admin@test.com / admin123

Ресурсы: products, orders, users

Структура
text
.
├── main.go      # точка входа, маршруты
├── auth.go      # JWT, bcrypt, права, БД
├── handlers.go  # обработчики запросов
├── models.go    # модели User, Resource, Permission
├── go.mod
└── Dockerfile
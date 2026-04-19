# Subscription Service

REST-сервис для агрегации данных об онлайн-подписках пользователей.

---

## 📋 Требования к заданию

- [x] CRUDL операции над подписками (Create, Read, Update, Delete, List).
- [x] Агрегация: подсчет суммарной стоимости подписок за период с фильтрацией.
- [x] Миграции для инициализации PostgreSQL.
- [x] Конфигурация через переменные окружения (.env).
- [x] Swagger-документация API.
- [x] Запуск сервиса через Docker Compose.

---

## 🛠 Технологический стек

- **Язык:** Go 1.22+
- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM/Driver:** pgx/v5
- **Config:** envconfig, godotenv
- **Logging:** uber-go/zap
- **Docs:** swaggo/gin-swagger
- **Migrations:** golang-migrate
- **Build:** Docker (Multi-stage build)

---

## ⚙️ Установка и запуск

Для запуска проекта вам потребуется установленный **Docker** и **Make**.

### 1. Клонирование репозитория

```bash
git clone https://github.com/chixxx1/subscription-service.git
cd subscription-service
```

### 2. Настройка конфигурации

Скопируйте пример файла окружения и настройте его при необходимости:

```bash
cp .env.example .env
```

Примечание: По умолчанию сервис настроен на подключение к PostgreSQL, запущенному в Docker Compose.


### 3. Запуск инфраструктуры
Поднимите базу данных и примените миграции:

```bash
# Запуск PostgreSQL
make env-up
# Применение миграций
make migrate-up
```

### 4. Запуск сервиса
Соберите образ приложения и запустите все сервисы:

```bash
make up
```

Сервис будет доступен по адресу: http://localhost:8080

---

📚 API Documentation 

После запуска сервиса документация Swagger доступна по адресу: 

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Основные эндпоинты 

| Метод | Путь | Описание |
| :--- | :--- | :--- |
| `POST` | `/api/v1/subscriptions` | Создать новую подписку |
| `GET` | `/api/v1/subscriptions` | Получить список подписок (с фильтрацией) |
| `GET` | `/api/v1/subscriptions/{id}` | Получить подписку по ID |
| `PUT` | `/api/v1/subscriptions/{id}` | Обновить подписку |
| `DELETE` | `/api/v1/subscriptions/{id}` | Удалить подписку |
| `GET` | `/api/v1/subscriptions/total-cost` | Подсчет общей стоимости за период |


Пример запроса на создание подписки 

```bash
POST /api/v1/subscriptions
```

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "04-2026"
}
```


🧪 Полезные команды Makefile 
 
| Команда | Описание |
| :--- | :--- |
| `make up` | Сборка и запуск всех сервисов (API + DB) |
| `make down` | Остановка всех контейнеров |
| `make env-up` | Запуск только базы данных |
| `make migrate-up` | Применение миграций БД |
| `make migrate-down` | Откат последней миграции |
| `make gen-swagger` | Перегенерация документации Swagger |
| `make start` | Локальный запуск Go-приложения (без Docker) |

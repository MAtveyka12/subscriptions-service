# Subscription Service

REST-сервис для агрегации данных об онлайн-подписках пользователей с возможностью CRUD операций и расчета стоимости подписок за выбранный период.

## 📋 Содержание

- [Описание проекта](#описание-проекта)
- [Технические требования](#технические-требования)
- [Архитектура](#архитектура)
- [API Документация](#api-документация)
- [Установка и запуск](#установка-и-запуск)
- [Тестирование](#тестирование)
- [Структура проекта](#структура-проекта)
- [Выполненные требования](#выполненные-требования)

## 🎯 Описание проекта

Subscription Service - это микросервис, предназначенный для управления подписками пользователей. Сервис предоставляет REST API для создания, чтения, обновления, удаления и поиска подписок, а также для расчета их стоимости за определенный период.

### Основные возможности:

- ✅ **CRUD операции** над подписками
- ✅ **Фильтрация** по пользователю и названию сервиса
- ✅ **Расчет стоимости** подписок за период
- ✅ **Пагинация** результатов
- ✅ **Валидация** входных данных
- ✅ **Логирование** всех операций
- ✅ **Swagger документация**

## 🛠 Технические требования

### Технологический стек:

- **Язык**: Go 1.24.3
- **Фреймворк**: Gin
- **База данных**: PostgreSQL 15+
- **Миграции**: Liquibase
- **Контейнеризация**: Docker & Docker Compose
- **Документация**: Swagger/OpenAPI 2.0
- **Логирование**: slog (структурированные логи)

### Системные требования:

- Docker 20.10+
- Docker Compose 2.0+
- 512MB RAM (минимум)
- 1GB свободного места

## 🏗 Архитектура

Проект следует принципам **Clean Architecture** с разделением на слои:

```
┌─────────────────────────────────────────┐
│              HTTP Layer                 │
│         (Controllers/Handlers)          │
├─────────────────────────────────────────┤
│           Application Layer             │
│              (Services)                 │
├─────────────────────────────────────────┤
│            Domain Layer                 │
│             (Models)                    │
├─────────────────────────────────────────┤
│         Infrastructure Layer            │
│         (Repository/Database)           │
└─────────────────────────────────────────┘
```

### Компоненты:

- **HTTP Handlers** - обработка HTTP запросов
- **Services** - бизнес-логика приложения
- **Repository** - работа с базой данных
- **Models** - доменные модели
- **DTO** - объекты передачи данных

## 📚 API Документация

### Базовый URL
```
http://localhost:8080
```

### Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/health` | Проверка состояния сервиса |
| `GET` | `/swagger/*` | Swagger документация |
| `POST` | `/subscriptions` | Создание подписки |
| `GET` | `/subscriptions` | Получение списка подписок |
| `GET` | `/subscriptions/{id}` | Получение подписки по ID |
| `PUT` | `/subscriptions/{id}` | Обновление подписки |
| `DELETE` | `/subscriptions/{id}` | Удаление подписки |
| `GET` | `/subscriptions/cost` | Расчет стоимости подписок |

### Модель данных

#### Subscription
```json
{
  "id": "uuid",
  "service_name": "string",
  "price": "integer",
  "user_id": "uuid",
  "start_date": "date",
  "end_date": "date (optional)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Примеры запросов

#### Создание подписки
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "2025-07-01",
    "end_date": "2025-12-31"
  }'
```

#### Расчет стоимости
```bash
curl -X GET "http://localhost:8080/subscriptions/cost?start_date=2025-07&end_date=2025-12&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

**Полная документация доступна по адресу:** `http://localhost:8080/swagger/index.html`

## 🚀 Установка и запуск

### Предварительные требования

1. Установите [Docker](https://docs.docker.com/get-docker/)
2. Установите [Docker Compose](https://docs.docker.com/compose/install/)

### ⚡ Быстрый старт (3 шага)

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd subscription-service
```

2. **Создайте файл окружения:**
```bash
cp .env.example .env
# Или вручную:
echo "DB_USER=db_user" > .env
echo "DB_PASSWORD=db_password" >> .env
```

3. **Запустите сервис:**
```bash
docker-compose up -d
```

4. **Проверьте статус:**
```bash
curl http://localhost:8080/health
```

### Остановка сервиса

```bash
docker-compose down
```

### Пересборка после изменений

```bash
docker-compose down
docker-compose up -d --build
```

## 🧪 Тестирование

### Автоматическое тестирование

Сервис включает health check эндпоинт для мониторинга:

```bash
curl http://localhost:8080/health
```

### Тестовые сценарии

1. **CRUD операции** - создание, чтение, обновление, удаление подписок
2. **Фильтрация** - поиск по пользователю и названию сервиса
3. **Расчет стоимости** - подсчет с различными фильтрами
4. **Валидация** - проверка обработки некорректных данных
5. **Пагинация** - работа с большими объемами данных

## 📁 Структура проекта

```
subscription-service/
├── cmd/
│   └── run/
│       └── main.go         # Точка входа приложения
├── config/
│   └── config.yaml         # Конфигурация сервиса
├── docs/
│   ├── docs.go             # Swagger документация
│   ├── swagger.json        # OpenAPI спецификация
│   └── swagger.yaml        # OpenAPI спецификация (YAML)
├── internal/
│   ├── app/
│   │   └── app.go          # Инициализация приложения
│   ├── application/
│   │   └── service/
│   │       ├── service.go  # Интерфейсы сервисов
│   │       └── subscription_service.go # Бизнес-логика подписок
│   ├── config/
│   │   └── config.go        # Загрузка конфигурации
│   ├── domain/
│   │   └── subscription/
│   │       └── models.go    # Доменные модели
│   └── infrastructure/
│       ├── controllers/
│       │   ├── dto/
│       │   │   ├── requests.go    # DTO для запросов
│       │   │   ├── responses.go   # DTO для ответов
│       │   │   └── date_parser.go # Кастомный парсер дат
│       │   └── http/
│       │       ├── handler.go     # Базовый обработчик
│       │       ├── handler_create.go      # Создание подписки
│       │       ├── handler_read.go        # Чтение подписки
│       │       ├── handler_update.go      # Обновление подписки
│       │       ├── handler_delete.go      # Удаление подписки
│       │       ├── handler_list.go        # Список подписок
│       │       ├── handler_calculate_cost.go # Расчет стоимости
│       │       ├── handler_helpers.go     # Вспомогательные функции
│       │       └── handler_register_routers.go # Регистрация маршрутов
│       └── repository/
│           └── subscription_repository.go # Работа с БД
├── migrations/
│   ├── 0001_create_subscription_table.sql # SQL миграция
│   └── master.xml               # Liquibase манифест
├── pkg/
│   └── http_server/
│       └── server.go            # HTTP сервер
├── docker-compose.yaml          # Docker Compose конфигурация
├── Dockerfile                   # Docker образ
├── go.mod                       # Go модули
├── go.sum                       # Go зависимости
├── .env.example                 # Пример файла окружения
└── README.md                    # Документация проекта
```

## ✅ Выполненные требования

### 1. CRUDL операции ✅

- **CREATE**: `POST /subscriptions` - создание подписки
- **READ**: `GET /subscriptions/{id}` - получение подписки по ID
- **UPDATE**: `PUT /subscriptions/{id}` - обновление подписки
- **DELETE**: `DELETE /subscriptions/{id}` - удаление подписки
- **LIST**: `GET /subscriptions` - получение списка подписок

### 2. Поля записи подписки ✅

- ✅ **service_name** - название сервиса (string)
- ✅ **price** - стоимость месячной подписки в рублях (integer)
- ✅ **user_id** - ID пользователя в формате UUID
- ✅ **start_date** - дата начала подписки (date)
- ✅ **end_date** - опциональная дата окончания подписки (date)

### 3. Расчет стоимости ✅

- ✅ **Эндпоинт**: `GET /subscriptions/cost`
- ✅ **Фильтрация по user_id** - UUID пользователя
- ✅ **Фильтрация по service_name** - название подписки
- ✅ **Период** - start_date и end_date в формате YYYY-MM
- ✅ **Сложная логика расчета** с учетом пересечений периодов

### 4. PostgreSQL с миграциями ✅

- ✅ **База данных**: PostgreSQL 15+
- ✅ **Миграции**: Liquibase
- ✅ **Автоматическое применение** при запуске
- ✅ **Индексы** для оптимизации запросов
- ✅ **UUID расширения** для генерации ID

### 5. Логирование ✅

- ✅ **Структурированные логи** через slog
- ✅ **Уровни логирования**: Debug, Info, Error
- ✅ **Контекстная информация** в логах
- ✅ **Логирование всех операций** в сервисном слое
- ✅ **Логирование HTTP запросов** через Gin middleware

### 6. Конфигурация ✅

- ✅ **YAML конфигурация** (`config/config.yaml`)
- ✅ **Поддержка .env файлов** для секретов
- ✅ **Вынесены все настройки**: порт, хост, БД, таймауты
- ✅ **Переменные окружения** для чувствительных данных

### 7. Swagger документация ✅

- ✅ **Полная документация** всех эндпоинтов
- ✅ **Интерактивный интерфейс** по адресу `/swagger/index.html`
- ✅ **Примеры запросов и ответов**
- ✅ **Валидация полей** в документации
- ✅ **OpenAPI 2.0 спецификация**

### 8. Docker Compose ✅

- ✅ **Полная конфигурация** для разработки и продакшена
- ✅ **PostgreSQL контейнер** с health checks
- ✅ **Приложение в контейнере** с multi-stage build
- ✅ **Автоматические миграции** через Liquibase
- ✅ **Сеть и volumes** для изоляции
- ✅ **Health checks** для всех сервисов

## 📊 Примеры использования

### Создание подписки на Yandex Plus
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "2025-07-01",
    "end_date": "2025-12-31"
  }'
```

### Расчет стоимости подписок пользователя за полгода
```bash
curl -X GET "http://localhost:8080/subscriptions/cost?start_date=2025-07&end_date=2025-12&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

### Поиск всех подписок на сервисы Yandex
```bash
curl -X GET "http://localhost:8080/subscriptions?service_name=Yandex"
```

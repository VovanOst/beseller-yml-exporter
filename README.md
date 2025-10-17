# BeSeller YML Exporter

Консольное приложение на Go для экспорта каталога товаров из BeSeller GraphQL API в формат YML (Yandex Market Language).

## Требования

- Go 1.22+
- Доступ к интернету

## Установка

```bash
# Клонировать репозиторий
git clone https://github.com/VovanOst/beseller-yml-exporter
cd beseller-yml-exporter

# Установить зависимости
go mod download

# Скопировать пример конфигурации
cp .env.example .env
```

## Конфигурация

Отредактируйте `.env` файл:

```env
GRAPHQL_ENDPOINT=https://demo.beseller.com/graphql?token=YOUR_TOKEN
SHOP_NAME=BeSeller Demo
SHOP_COMPANY=ООО Открытый контакт
SHOP_URL=https://demo.beseller.com
CURRENCY=BYN
STATUS_ID=1
OUTPUT_PATH=export.yml
HTTP_TIMEOUT=30
LOG_LEVEL=info
```

## Запуск

```bash
# Использование переменных из .env
make run

# Или напрямую через go
go run cmd/exporter/main.go

# С переопределением параметров через флаги
go run cmd/exporter/main.go \
  --endpoint="https://demo.beseller.com/graphql?token=YOUR_TOKEN" \
  --out=export.yml \
  --shop-name="My Shop" \
  --shop-company="My Company" \
  --shop-url="https://myshop.com" \
  --currency=BYN \
  --status-id=1 \
  --timeout=30s \
  --log-level=debug
```

## Сборка

```bash
# Сборка бинарного файла
make build

# Запуск бинарника
./bin/exporter
```

## Архитектура

Проект следует принципам Clean Architecture:

```
cmd/exporter/          - Точка входа приложения
internal/
  domain/              - Бизнес-логика и интерфейсы
    entity/            - Доменные сущности
    repository/        - Интерфейсы репозиториев
  usecase/             - Сценарии использования
  infrastructure/      - Технические детали
    graphql/           - GraphQL клиент и репозиторий
    yml/               - YML writer
    config/            - Конфигурация
  logger/              - Логирование
pkg/
  errors/              - Кастомные ошибки
```

## Описание работы

1. Загружает конфигурацию из `.env` и флагов командной строки
2. Подключается к GraphQL API BeSeller
3. Получает список всех категорий
4. Получает товары со статусом "новинка" (statusId=1)
5. Формирует YML файл в формате Yandex Market Language
6. Сохраняет результат в указанный файл

## Формат YML

Приложение генерирует YML файл согласно спецификации Яндекс.Маркет:
- Категории с поддержкой иерархии (parentId)
- Товары (offers) с полями: url, price, currency, category, pictures, name, vendor, barcode, description
- Только товары со статусом "новинка" (statusId=1)
- Атрибут available на основе данных о наличии

## Makefile команды

```bash
make help        # Показать доступные команды
make run         # Запустить приложение
make build       # Собрать бинарный файл
make test        # Запустить тесты
make clean       # Очистить бинарные файлы
make fmt         # Форматировать код
make lint        # Проверить код линтером
```

## Логирование

Уровни логирования: `debug`, `info`, `warn`, `error`

Установить через флаг `--log-level` или переменную `LOG_LEVEL` в `.env`

## Обработка ошибок

- HTTP ошибки: автоматический retry с exponential backoff (3 попытки)
- GraphQL errors: логирование и завершение с кодом ошибки
- Валидация данных: пропуск некорректных товаров с предупреждением
- Частичные данные: продолжение работы с логированием

## Ограничения

- Поддерживается только GraphQL API BeSeller
- Экспортируются только товары со статусом "новинка" (statusId=1)
- Формат YML соответствует базовой спецификации Yandex Market Language
- Опциональные поля (vendor, barcode) пропускаются если отсутствуют

## Пример вывода

```
2025-10-15 14:30:00 INFO Starting BeSeller YML Exporter
2025-10-15 14:30:00 INFO Connecting to GraphQL endpoint
2025-10-15 14:30:01 INFO Fetching categories...
2025-10-15 14:30:02 INFO Found 45 categories
2025-10-15 14:30:02 INFO Fetching products with statusId=1...
2025-10-15 14:30:04 INFO Found 128 products
2025-10-15 14:30:04 INFO Generating YML file...
2025-10-15 14:30:04 INFO YML file created: export.yml
2025-10-15 14:30:04 INFO Export completed successfully (categories=45, offers=128)
```

## Разработка

### Добавление новых источников данных

Благодаря чистой архитектуре можно легко добавить альтернативные источники:

1. Реализовать интерфейс `domain.CatalogRepository`
2. Добавить новый адаптер в `internal/infrastructure/`
3. Инжектировать в `cmd/exporter/main.go`

### Добавление новых форматов экспорта

1. Реализовать интерфейс `usecase.CatalogWriter`
2. Добавить writer в `internal/infrastructure/`
3. Использовать в use case

## Лицензия

MIT

## Автор

Владимир Островский

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## Локальный запуск

Поднять PostgreSQL и MongoDB:

```bash
task up-all
```

Создать локальные конфигурации сервисов:

```bash
cp order/.env.example order/.env
cp inventory/.env.example inventory/.env
```

Запустить `inventory`, `payment` и `order` из директорий соответствующих модулей, затем проверить весь API-сценарий:

```bash
task test-api
```

Остановить зависимости и удалить локальные volumes:

```bash
task down-all
```

Миграции Order Service встроены в бинарник через `go:embed` и автоматически применяются при старте сервиса.

## Проверки

```bash
task test
task lint
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Unit-тесты
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml

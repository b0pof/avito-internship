# Тестовое задание на позицию стажера Avito

Условие и требования можно найти в папке `/задание`.  
API реализовано на языке **Go**.  
СУБД: **Postgres**.  
Миграции: **goose**.  
Старался писать чистый код и придерживаться чистой архитектуры.

## Функционал

Весь основной функционал реализован.  
Из **дополнительных** требований:  
1. Версионирование тендеров и предложений, возможность редактирования и отката версии;
2. Описание конфигурации линтера (`golangci.yml`).

API приложения описано в `/postman`.

## Запуск приложения

### Конфигурация

Конфигурация приложения производится через следующие переменные окружения:
- `SERVER_ADDRESS` — адрес и порт, который будет слушать HTTP сервер при запуске. Пример: 0.0.0.0:8080.
- `POSTGRES_CONN` — URL-строка для подключения к PostgreSQL в формате postgres://{username}:{password}@{host}:{5432}/{dbname}.
- `POSTGRES_JDBC_URL` — JDBC-строка для подключения к PostgreSQL в формате jdbc:postgresql://{host}:{port}/{dbname}.
- `POSTGRES_USERNAME` — имя пользователя для подключения к PostgreSQL.
- `POSTGRES_PASSWORD` — пароль для подключения к PostgreSQL.
- `POSTGRES_HOST` — хост для подключения к PostgreSQL (например, localhost).
- `POSTGRES_PORT` — порт для подключения к PostgreSQL (например, 5432).
- `POSTGRES_DATABASE` — имя базы данных PostgreSQL, которую будет использовать приложение.

### Команды для запуска

В корне проекта можно найти `Makefile`.  
Команда `make` без аргументов выводит возможные команды и их описание.  
Минимальное множество команд для запуска приложения:
1. `make run-db` (поднимает Docker-контейнер с Postgres на `localhost:5432`)
2. `make migrations-up` (если схема данных не создана)
3. `make run-app` (если адрес сервера не указан, приложение запускается на `0.0.0.0:8080`)

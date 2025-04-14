# Сервис для работы с ПВЗ

## Стек
Язык сервиса: Go.
База данных: PostgreSQL.

## Запуск

docker-compose -f deploy/compose.yml up

## Список необходимых переменных окружения:
Переменные окружения попадают в контейнер из .env файла, который нужно расположить в корневой директории
* APPLICATION_PORT=8080
* SECRET_KEY=my_secret_key
* POSTGRES_USER=admin
* POSTGRES_PASSWORD=123
* POSTGRES_DB=pvz
* DATABASE_HOST=localhost
## Доступ по порту 8080, у всех апи ручек добавлен префикс /api/v1, пример запроса:
  localhost:8080/api/v1/dummyLogin

## Для доступа к апи с необходимостью авторизации необходим http заголовок Authorization:
  Bearer {Токен}
  ![image](https://github.com/user-attachments/assets/552e0c5d-ab7f-47a6-874b-e67ab2441b67)


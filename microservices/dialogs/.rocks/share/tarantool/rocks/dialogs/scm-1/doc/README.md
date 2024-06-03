# Dialogs Microservice

## Запуск
1. установить tarantool 3 https://www.tarantool.io/en/doc/latest/how-to/getting_started_db/
2. установить tarantool cli (tt) https://www.tarantool.io/en/doc/latest/how-to/getting_started_db/
3. make build
4. make start
5. делаем http запросы на http://localhost:8083/api/get_dialog (GET) и http://localhost:8083/api/send_message (POST)
список параметров смотрите в router.lua

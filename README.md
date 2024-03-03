# Highload Social Network (hlsoc)

Мой учебный проект в рамках курса Highload Architect от Otus (https://otus.ru/lessons/highloadarchitect/).

Так как учеба в целом не о коде - я не использовал лучшие парктики и паттерны разработки.

## Запуск

1. Создать файл .env на оснвое env.example
2. Запустить проект с помощью `make docker-run`

## Взаимодействие

Через GRPC. Рефлексия включена - сервер делится с клиентом protobuf схемой.
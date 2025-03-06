# Garantex API Go GRPC Example

Сервер который отправляет HTTP-запрос к garantex.org (https://garantex.org/api/v2/depth) для получения курса USDT по всем текущим рынкам и возвращает ответ по gRPC. Вместе с сервером тестовый gRPC клиент

# Запуск

Приложение вместе с клиентом и зависимостями запускается в docker compose с помощью команды make docker build, контейнер останавливается командой make docker-stop

```bash
make docker-build
make docker-stop
```
Также сервер и клиент можно запустить отдельно для передачи флагов командной строки  

```bash
make build //компиляция сервера
docker compose up //создание образа с зависимостями
make run //запуск сервера
go run ./cmd - db-max-open-conns=100 - db-max-idle-conns=50 //запуск с пользовательскими флагами
go run ./client  //запуск тестового клиента
```

Для запуска автотестов используйте make test, для проверки линтером make lint



## License

[MIT](https://choosealicense.com/licenses/mit/)
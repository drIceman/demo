## Демонстрационный проект

### Инструмент сборки образов для Go [https://ko.build]

```
curl -sSfL https://github.com/ko-build/ko/releases/download/v0.15.1/ko_Linux_x86_64.tar.gz | tar xz ko
chmod +x ./ko
sudo mv ./ko /usr/bin
```

### Использование

Склонировать репозиторий и перейти в директорию проекта

#### Приложение demo

Сборка образа

`ko build -PL ./cmd/demo`

Запуск приложения

`docker run -it --rm ko.local/github.com/driceman/demo/cmd/demo:latest`

#### Приложение api-server

По мотивам [https://go.dev/doc/tutorial/web-service-gin], но с книгами

Сборка образа

`ko build -PL ./cmd/api-server`

Запуск приложения

`docker run -it --rm -p 8080:8080 ko.local/github.com/driceman/demo/cmd/api-server:latest`

можно добавить к вызову `--env GIN_MODE=release`

Получение списка книг

`curl 127.0.0.1:8080/books`

Получение книги по ID

`curl 127.0.0.1:8080/books/2`

Создание книги

```
curl 127.0.0.1:8080/books \
    -H "Content-Type: application/json" \
    -X POST \
    --data '{"id": "4", "title": "Книга 4", "author": "Автор 4", "price": 149.99}'
```

#### Приложение files-parser

Сборка образа

`ko build -PL ./cmd/files-parser`

Запуск приложения

`docker run -it --rm -v ./internal/csv-parser/stub.csv:/app/stub.csv ko.local/github.com/driceman/demo/cmd/files-parser:latest --file_type csv --file_path /app/stub.csv --from_byte 0 --rows_limit 11 --threads_count 1`

или

`docker run -it --rm -v ./internal/csv-parser/stub.csv:/app/stub.csv -e FILE_TYPE csv -e FILE_PATH /app/stub.csv -e FROM_BYTE 0 -e ROWS_LIMIT 11 -e THREADS_COUNT 1 ko.local/github.com/driceman/demo/cmd/files-parser:latest`


Тестирование

`go test -timeout 30s -run=. github.com/drIceman/demo/internal/csv-parser`

Бенчмарк

`go test -benchmem -run=. -bench=. github.com/drIceman/demo/internal/csv-parser`

Тестирование с файлом > 1Гб

создать большой файл
`TEST_1G=1 go test -run=^TestPrepare1G$ github.com/drIceman/demo/internal/csv-parser`

запустить тест
`TEST_1G=1 go test -v -benchmem -run=^BenchmarkParse1G$ -bench=^BenchmarkParse1G$ github.com/drIceman/demo/internal/csv-parser`

Профилирование (на маленьких объемах не срабатывает)

```
go build ./cmd/files-parser
./files-parser --file_type csv --file_path ./internal/csv-parser/stub.csv --from_byte 0 --rows_limit 11 --threads_count 1 --mem_profile_path prof.mprof
go tool pprof files-parser prof.mprof
```

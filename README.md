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

`docker run -it --rm ko.local/demo/cmd/demo:latest`

#### Приложение api-server

По мотивам [https://go.dev/doc/tutorial/web-service-gin], но с книгами

Сборка образа

`ko build -PL ./cmd/api-server`

Запуск приложения

`docker run -it --rm -p 8080:8080 ko.local/demo/cmd/api-server:latest`

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


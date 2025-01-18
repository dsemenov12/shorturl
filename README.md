# Сервис сокращения URL

## Описание

Этот проект представляет собой сервис для сокращения длинных URL-адресов, позволяя преобразовывать их в короткие ссылки для удобного использования и обмена.

## Функциональные возможности

- **Сокращение URL**: преобразование длинных ссылок в короткие.
- **Перенаправление**: переход по короткой ссылке ведет на исходный длинный URL.
- **Статистика**: отслеживание количества переходов по каждой короткой ссылке.
- **Удаление ссылок**: возможность удаления ранее созданных коротких ссылок.
- **Обновление ссылок**: изменение целевого URL для существующей короткой ссылки.
- **Получение списка ссылок**: возможность получить список всех созданных коротких ссылок.

## Установка и запуск

1. **Клонирование репозитория**:

   ```bash
   git clone https://github.com/dsemenov12/shorturl.git
   ```

2. **Инициализация модуля**:

   ```bash
   go mod init github.com/dsemenov12/shorturl
   ```

3. **Установка зависимостей**:

   ```bash
   go mod tidy
   ```

4. **Запуск сервиса**:

   ```bash
   go run cmd/shorturl/main.go
   ```

   По умолчанию сервис будет доступен по адресу `http://localhost:8080`.

## Методы API

### 1. Сокращение URL

**POST** `/api/shorten`  

#### Запрос:
```json
{
  "url": "https://example.com/long-url"
}
```

#### Ответ:
```json
{
  "short_url": "http://localhost:8080/abcd123"
}
```

---

### 2. Переход по короткой ссылке

**GET** `/{short_url}`  

- Перенаправляет пользователя на исходный длинный URL.

---

### 3. Получение статистики по ссылке

**GET** `/api/stats/{short_url}`  

#### Ответ:
```json
{
  "original_url": "https://example.com/long-url",
  "clicks": 25
}
```

---

### 4. Удаление короткой ссылки

**DELETE** `/api/delete/{short_url}`  

#### Ответ:
```json
{
  "message": "Short URL deleted successfully"
}
```

---

### 5. Обновление длинного URL

**PUT** `/api/update/{short_url}`  

#### Запрос:
```json
{
  "new_url": "https://new-example.com/updated-url"
}
```

#### Ответ:
```json
{
  "message": "Short URL updated successfully"
}
```

---

### 6. Получение всех ссылок

**GET** `/api/links`  

#### Ответ:
```json
[
  {
    "short_url": "http://localhost:8080/abcd123",
    "original_url": "https://example.com/long-url",
    "clicks": 25
  },
  {
    "short_url": "http://localhost:8080/xyz789",
    "original_url": "https://another-example.com",
    "clicks": 10
  }
]
```

## Тестирование

Для запуска тестов выполните:

```bash
go test ./...
```


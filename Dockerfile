# Используем официальный образ Go как базовый образ
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum (если они есть) и загружаем зависимости
COPY go.mod go.sum .env ./
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o go_final_project .

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранный бинарный файл из builder-стадии
COPY --from=builder /app/go_final_project .
COPY --from=builder /app/.env .

# Открываем порт, который будет использовать приложение
EXPOSE 7540

# Команда для запуска приложения
CMD ["./go_final_project"]
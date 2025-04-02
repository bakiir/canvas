# Canvas Application

Canvas Application – это учебный проект для изучения разработки веб-приложений на Go с использованием GORM, Gin и PostgreSQL в Docker.

## 🚀 Функционал
- Управление студентами, курсами и преподавателями
- Авторизация и аутентификация
- Работа с базой данных PostgreSQL через GORM
- REST API с использованием Gin
- Docker-контейнеризация

## 🛠️ Технологии
- Go (Gin, GORM)
- PostgreSQL
- Docker
- Swagger (для документации API)

## 📦 Установка и запуск
### 1️⃣ Клонирование репозитория
```sh
git clone https://github.com/bakiir/canvas.git
cd canvas
```

### 2️⃣ Запуск в Docker
```sh
docker-compose up --build
```

### 3️⃣ Запуск вручную
1. Установите PostgreSQL и создайте базу данных `canvas_db`
2. Настройте переменные окружения (`.env`)
3. Запустите сервер
   ```sh
   go run main.go
   ```

## 🔗 API Документация
После запуска сервер доступен по адресу: `http://localhost:8080`  

## 👤 Авторы
- [bakiir](https://github.com/bakiir)



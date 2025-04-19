DB_URL=postgres://postgres:postgres@localhost:5439/canvas_db?sslmode=disable
MIGRATIONS_DIR=./migrations
NAME=my_migration
VERSION=1

.PHONY: up down down-all force goto version create install-migrate migrate-force

# Применить все миграции
up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Откатить последнюю миграцию
down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Откатить ВСЕ миграции
down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

# Принудительно задать версию миграции (указать VERSION=номер)
force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

# Перейти к конкретной версии миграции (указать VERSION=номер)
goto:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" goto $(VERSION)

# Посмотреть текущую версию миграций
version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Создать новую миграцию (указать NAME=имя_миграции)
create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Используй migrate force, чтобы сбросить dirty-флаг
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

# Установить миграции
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

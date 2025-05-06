#!/bin/sh

echo "Running database migration..."
./migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "Migration complete. Starting service..."
./canvas-service

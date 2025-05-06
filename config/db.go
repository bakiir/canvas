package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// Connect инициализирует подключение к базе данных// locLHOST wrap canvas_db
func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment")
	}
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	sslmode := "disable"

	// Формируем строку подключения правильно
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbPort, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.Exec("SET session_replication_role = 'origin'")

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)                  // Максимальное количество бездействующих соединений
	sqlDB.SetMaxOpenConns(100)                 // Максимальное количество открытых соединений
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Максимальное время жизни соединения

	log.Println("🔹 Using external migration tool (golang-migrate). No AutoMigrate is run.")

	DB = db
	log.Println("Successfully connected to database")

}

type StudentCourse struct {
	gorm.Model
	StudentID uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

package config

import (
	"CanvasApplication/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

// Connect инициализирует подключение к базе данных
func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5439/postgres?sslmode=disable"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Конфигурация пула соединений
	sqlDB.SetMaxIdleConns(10)                  // Максимальное количество бездействующих соединений
	sqlDB.SetMaxOpenConns(100)                 // Максимальное количество открытых соединений
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Максимальное время жизни соединения в минутах

	// Автомиграция моделей
	if err := migrateModels(db); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	DB = db
	log.Println("Successfully connected to database")
}

func migrateModels(db *gorm.DB) error {
	// Отключаем проверки внешних ключей
	if err := db.Exec("SET session_replication_role = replica").Error; err != nil {
		return fmt.Errorf("failed to disable constraints: %w", err)
	}

	//// Удаляем таблицы
	//db.Migrator().DropTable(&StudentCourse{})
	//db.Migrator().DropTable(&models.Course{})
	//db.Migrator().DropTable(&models.Student{})
	//db.Migrator().DropTable(&models.Teacher{})

	// Сбрасываем последовательности (важно!)
	db.Exec("ALTER SEQUENCE students_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE teachers_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE courses_id_seq RESTART WITH 1")

	// Создаем таблицы
	if err := db.AutoMigrate(&models.Teacher{}, &models.Student{}, &models.Course{}); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	// Включаем проверки внешних ключей обратно
	if err := db.Exec("SET session_replication_role = DEFAULT").Error; err != nil {
		return fmt.Errorf("failed to enable constraints: %w", err)
	}

	return nil
}

type StudentCourse struct {
	gorm.Model
	StudentID uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

// GetDB возвращает экземпляр базы данных
func GetDB() *gorm.DB {
	if DB == nil {
		panic("Database connection not initialized. Call Connect() first")
	}
	return DB
}

// CloseDB закрывает соединение с базой данных
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Failed to close database connection: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}

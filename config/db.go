package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

// Connect инициализирует подключение к базе данных
func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5439/canvas_db?sslmode=disable"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.Exec("SET session_replication_role = 'origin'")

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
	log.Println("🔹 Starting migration...")

	// Check which database you're connected to
	dbName := ""
	db.Raw("SELECT current_database()").Scan(&dbName)
	log.Printf("🔹 Connected to database: %s", dbName)

	// Set replication role for disabling foreign key checks
	db.Exec("SET session_replication_role = 'replica'")

	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = 'origin'")
	log.Println("✅ Database migration completed successfully.")
	return nil
}

type StudentCourse struct {
	gorm.Model
	student_id uint `gorm:"primaryKey"`
	course_id  uint `gorm:"primaryKey"`
	CreatedAt  time.Time
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

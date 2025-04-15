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
	log.Println("🔹 Using external migration tool (golang-migrate). No AutoMigrate is run.")

	DB = db
	log.Println("Successfully connected to database")
}

//func migrateModels(db *gorm.DB) error {
//	log.Println("🔹 Starting migration...")
//
//	// Check which database you're connected to
//	dbName := ""
//	db.Raw("SELECT current_database()").Scan(&dbName)
//	log.Printf("🔹 Connected to database: %s", dbName)
//
//	// Set replication role for disabling foreign key checks
//	db.Exec("SET session_replication_role = 'replica'")
//
//	log.Println("🔹 Dropping existing tables...")
//	if err := db.Migrator().DropTable(&models.Teacher{}, &models.Student{}, &models.Course{}, &models.StudentCourse{}); err != nil {
//		log.Printf("⚠️ Error dropping tables: %v", err)
//		return fmt.Errorf("failed to drop tables: %w", err)
//	}
//
//	log.Println("🔹 Running AutoMigrate...")
//	if err := db.AutoMigrate(&models.Teacher{}, &models.Student{}, &models.Course{}, &models.StudentCourse{}); err != nil {
//		log.Printf("⚠️ AutoMigrate error: %v", err)
//		return fmt.Errorf("failed to migrate models: %w", err)
//	}
//
//	// Re-enable foreign key checks
//	db.Exec("SET session_replication_role = 'origin'")
//	log.Println("✅ Database migration completed successfully.")
//	return nil
//}

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

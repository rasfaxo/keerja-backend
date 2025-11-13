package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger
	var gormLogger gormlogger.Interface
	if cfg.IsDevelopment() {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Info)
	} else {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Error)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Prepared statement cache
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	DB = db
	return db, nil
}

// CloseDB closes the database connection gracefully
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed successfully")
	return nil
}

// AutoMigrate runs auto migration for all models
func AutoMigrate(db *gorm.DB) error {
	// Import all domain models here when they're created
	// Example:
	// err := db.AutoMigrate(
	// 	&user.User{},
	// 	&user.UserProfile{},
	// 	&company.Company{},
	// 	&job.Job{},
	// 	// ... other models
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to run auto migration: %w", err)
	// }

	log.Println("Auto migration completed (Note: Add models when they're created)")
	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDB() first")
	}
	return DB
}

// InitDatabase is an alias for InitDB that loads config and initializes database
func InitDatabase() error {
	cfg := LoadConfig()
	db, err := InitDB(cfg)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// CloseDatabase is an alias for CloseDB
func CloseDatabase() error {
	return CloseDB()
}

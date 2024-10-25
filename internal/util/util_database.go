package util

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	POSTGRES = "postgres"
	MYSQL    = "mysql"
)

func NewLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			SlowThreshold:             200 * time.Millisecond,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	return newLogger
}

func NewPostgresDB() *gorm.DB {
	dsn := "host=" + config.DB_POSTGRES.DB_HOST + " user=" + config.DB_POSTGRES.DB_USER + " password=" + config.DB_POSTGRES.DB_PASSWORD + " dbname=" + config.DB_POSTGRES.DB_NAME + " port=" + config.DB_POSTGRES.DB_PORT + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewLogger(),
	})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil
	}
	return db
}

func NewMySQLDB() *gorm.DB {
	dsn := config.DB_MYSQL.DB_USER + ":" + config.DB_MYSQL.DB_PASSWORD + "@tcp(" + config.DB_MYSQL.DB_HOST + ":" + config.DB_MYSQL.DB_PORT + ")/" + config.DB_MYSQL.DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewLogger(),
	})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil
	}
	return db
}

func SetupDatabaseConnection(SGBD string) (*gorm.DB, *sql.DB, error) {
	var db *gorm.DB

	switch SGBD {
	case POSTGRES:
		db = NewPostgresDB()
	case MYSQL:
		db = NewMySQLDB()
	default:
		return nil, nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, sqlDB, nil
}

func CheckConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

func Shutdown(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
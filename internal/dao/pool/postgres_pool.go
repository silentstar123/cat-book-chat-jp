package pool

import (
	"fmt"
	"log"
	"os"
	"time"

	"chat-room/config"
	"chat-room/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresDB *gorm.DB

func InitPostgresDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	pgConfig := config.GetConfig().Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		pgConfig.Host,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.DBName,
		pgConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get postgres database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库结构
	err = db.AutoMigrate(&model.PostgresUser{}, &model.UserProfiles{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	postgresDB = db
}

func GetPostgresDB() *gorm.DB {
	return postgresDB
} 
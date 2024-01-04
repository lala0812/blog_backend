package controller

import (
	"fmt"
	"go_gin_example/envconfig"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	// dsn 資料庫的連線資訊
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Taipei", envconfig.GetEnv("DB_HOST"), envconfig.GetEnv("DB_USER"), envconfig.GetEnv("DB_PASSWORD"), envconfig.GetEnv("DB_NAME"), envconfig.GetEnv("DB_PORT"), envconfig.GetEnv("DB_WITH_SSL"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to close database")
	}
	sqlDB.Close()
}


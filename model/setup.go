package model

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host= " + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_DBNAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable " +
		" timezone=Asia/Shanghai"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get DB instance: ", err)
	}

	// Batasi jumlah koneksi
	sqlDB.SetMaxOpenConns(5) // maksimal koneksi aktif
	sqlDB.SetMaxIdleConns(2) // idle connection
	sqlDB.SetConnMaxLifetime(time.Hour)

	// DB = database.Debug()
	DB = database
}

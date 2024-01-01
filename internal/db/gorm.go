package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConn() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost", "postgres", "postgres",
		"crypto-sig", "5432",
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(fmt.Errorf("не удалось подключиться к бд %s", err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(250)
	return db
}

func NewDBConnHost(host string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, "postgres", "postgres",
		"crypto-sig", "5432",
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(fmt.Errorf("не удалось подключиться к бд %s", err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(250)
	return db
}

// func Migrate(db *gorm.DB) {
// 	db.AutoMigrate(&models.CryptoPrice{})
// 	db.AutoMigrate(&models.UserClosedOrder{})
// 	db.AutoMigrate(&models.UserOpenOrder{})
// 	db.AutoMigrate(&models.KcsOrderModel{})
// }

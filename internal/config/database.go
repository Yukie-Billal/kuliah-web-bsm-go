package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	driver := os.Getenv("DB_DRIVER")

	switch driver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case "mysql", "":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	default:
		panic("unsupported DB_DRIVER: " + driver)
	}

	if err != nil {
		panic(err)
	}

	return db
}

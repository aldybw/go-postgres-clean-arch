package config

import (
	"fmt"
	"go-postgres-clean-arch/helper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbName   = "clean_arch_test"
)

func DatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	helper.ErrorPanic(err)

	return db
}

package database

import (
	"github.com/seasea128/WebAPI/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitConnection(cfg *config.Configuration) (*gorm.DB, error) {
	var connection gorm.Dialector

	switch cfg.DBType {
	case config.Postgres:
		connection = postgres.Open(cfg.ConnectionString)
	case config.SQLite:
		connection = sqlite.Open(cfg.ConnectionString)
	}

	db, err := gorm.Open(connection)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate()
	return db, nil
}

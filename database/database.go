package database

import (
	"github.com/seasea128/WebAPI/config"
	"github.com/seasea128/WebAPI/database/migrations"
	"github.com/seasea128/WebAPI/database/model"
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

	// TODO: Implement migration system
	db.AutoMigrate(model.Sessions{}, model.SuspensionLogs{}, model.Users{})
	return db, nil
}

func loadMigrations(db *gorm.DB, migrations []migrations.Migrations) {

}

package model

import (
	"database/sql"
	"time"
)

type Sessions struct {
	ID           uint64
	Name         string `gorm:"unique"`
	ControllerID string
	SessionID    int32
	CreatedAt    time.Time
	FinishedAt   sql.NullTime
	DeletedAt    sql.NullTime
}

package model

import (
	"database/sql"
	"time"
)

type Sessions struct {
	ID           uint64
	Name         string
	ControllerID string
	SessionID    int32
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
}

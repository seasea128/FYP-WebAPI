package model

import (
	"database/sql"
	"time"
)

type Sessions struct {
	ID        uint64
	Name      string
	UserID    uint64
	User      Users
	CreatedAt time.Time
	DeletedAt sql.NullTime
}

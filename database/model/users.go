package model

import "time"

type Users struct {
	ID             uint64
	CreatedAt      time.Time
	Username       string
	HashedPassword string
	UpdatedAt      time.Time
}

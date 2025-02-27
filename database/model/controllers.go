package model

import "time"

type Controllers struct {
	ID             uint64
	CreatedAt      time.Time
	ControllerName string
	UpdatedAt      time.Time
}

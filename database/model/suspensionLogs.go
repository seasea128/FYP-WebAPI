package model

import "time"

type SuspensionLogs struct {
	ID          uint64
	CreatedAt   time.Time
	SessionID   uint64
	Session     Sessions
	LeftTop     string
	LeftBottom  string
	RightTop    string
	RightBottom string
	GPSPosition string
	GPSSpeed    string
}

package model

import "time"

type SuspensionLogs struct {
	ID           uint64
	CreatedAt    time.Time
	ControllerID string
	SessionID    int32
	Session      Sessions
	LeftTop      int32
	LeftBottom   int32
	RightTop     int32
	RightBottom  int32
	GPSPosition  string
	GPSSpeed     string
}

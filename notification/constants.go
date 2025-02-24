package notification

import (
	"time"
)

type MediaServer int

const (
	Emby MediaServer = iota
)

type Metadata struct {
	Server   MediaServer
	Username string
	ImdbId   string
	Time     time.Time
}

type WatchedNotification struct {
	Metadata
	Watched bool
}

type PlaybackNotification struct {
	Metadata
	Playing  bool
	Position time.Duration
	Runtime  time.Duration
}

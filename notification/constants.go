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

type PlaybackAction int

const (
	Start PlaybackAction = iota
	Play
	Pause
	Stop
)

type PlaybackNotification struct {
	Metadata
	Action   PlaybackAction
	Position time.Duration
	Runtime  time.Duration
}

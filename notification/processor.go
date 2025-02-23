package notification

import (
	"fmt"
	"log/slog"
)

import "github.com/computer-geek64/emboxd/letterboxd"

type Processor struct {
	callback func(letterboxd.Event)
}

func NewProcessor(callback func(letterboxd.Event)) Processor {
	return Processor{
		callback: callback,
	}
}

func (n Processor) ProcessWatchedNotification(notification WatchedNotification) {
	var action letterboxd.Action
	if notification.Watched {
		action = letterboxd.FilmWatched
	} else {
		action = letterboxd.FilmUnwatched
	}

	n.callback(letterboxd.Event{
		ImdbId: notification.ImdbId,
		Action: action,
		Time:   notification.Time,
	})

	slog.Info(fmt.Sprintf("Processing watched notification %+v", notification))
}

func (n Processor) ProcessPlaybackNotification(notification PlaybackNotification) {
	// TODO: setup DB
	slog.Warn(fmt.Sprintf("Unprocessed playback notification %+v", notification))
}

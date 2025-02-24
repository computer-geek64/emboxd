package notification

import (
	"fmt"
	"log/slog"
	"time"
)

import "github.com/computer-geek64/emboxd/letterboxd"

// Watched percentage of total runtime to log movie as actually watched
const _MIN_WATCHED_PERCENTAGE int = 80

type Processor struct {
	callback                          func(letterboxd.Event)
	watchedDurationByImdbId           map[string]time.Duration
	playbackStartNotificationByImdbId map[string]PlaybackNotification
}

func NewProcessor(callback func(letterboxd.Event)) Processor {
	return Processor{
		callback:                          callback,
		watchedDurationByImdbId:           make(map[string]time.Duration),
		playbackStartNotificationByImdbId: make(map[string]PlaybackNotification),
	}
}

func (n *Processor) ProcessWatchedNotification(notification WatchedNotification) {
	slog.Info(fmt.Sprintf("Processing watched notification %+v", notification))

	var action letterboxd.Action
	if notification.Watched {
		action = letterboxd.FilmWatched
	} else {
		action = letterboxd.FilmUnwatched
	}

	delete(n.watchedDurationByImdbId, notification.ImdbId)
	delete(n.playbackStartNotificationByImdbId, notification.ImdbId)

	n.callback(letterboxd.Event{
		ImdbId: notification.ImdbId,
		Action: action,
		Time:   notification.Time,
	})
}

func (n *Processor) ProcessPlaybackNotification(notification PlaybackNotification) {
	slog.Info(fmt.Sprintf("Processing playback notification %+v", notification))

	// TODO: setup DB for permanent storage of partially watched films
	if notification.Playing {
		if _, alreadyStarted := n.playbackStartNotificationByImdbId[notification.ImdbId]; !alreadyStarted {
			// Keep earliest playback notification for current session
			n.playbackStartNotificationByImdbId[notification.ImdbId] = notification
		}
	} else {
		if startNotification, hasStart := n.playbackStartNotificationByImdbId[notification.ImdbId]; hasStart {
			var watchedDuration = min(
				// Ensure movie was actually watched
				notification.Time.Sub(startNotification.Time),
				// Ensure rewinding/replaying is not included in watched duration
				max(notification.Position-startNotification.Position, 0),
			)

			if _, partiallyWatched := n.watchedDurationByImdbId[notification.ImdbId]; partiallyWatched {
				n.watchedDurationByImdbId[notification.ImdbId] += watchedDuration
			} else {
				n.watchedDurationByImdbId[notification.ImdbId] = watchedDuration
			}
			delete(n.playbackStartNotificationByImdbId, notification.ImdbId)
		} else {
			n.watchedDurationByImdbId[notification.ImdbId] = notification.Position
			slog.Warn("Missing playback start time, set total watched duration to current playback position")
		}

		var watchedPercentage = float64(n.watchedDurationByImdbId[notification.ImdbId].Nanoseconds()) / float64(notification.Runtime.Nanoseconds()) * 100
		if int(watchedPercentage) >= _MIN_WATCHED_PERCENTAGE {
			delete(n.watchedDurationByImdbId, notification.ImdbId)

			n.callback(letterboxd.Event{
				ImdbId: notification.ImdbId,
				Action: letterboxd.FilmLogged,
				Time:   notification.Time,
			})
		}
	}
}

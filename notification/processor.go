package notification

import (
	"fmt"
	"log/slog"
	"time"
)

import "github.com/computer-geek64/emboxd/letterboxd"

// Watched percentage of total runtime to log movie as actually watched
const _MIN_WATCHED_PERCENTAGE uint = 70

// Position percentage of total runtime to evaluate if movie should be logged
const _MIN_POSITION_PERCENTAGE uint = 90

// Max elapsed time to ignore consecutive playback stop notifications
const _MAX_DUPLICATE_STOP_PLAYBACK_ELAPSED_TIME time.Duration = 2 * time.Minute

type Processor struct {
	callback                          func(letterboxd.Event)
	watchedDurationByImdbId           map[string]time.Duration
	playbackStartNotificationByImdbId map[string]PlaybackNotification
	playbackStopTimeByImdbId          map[string]time.Time
}

func NewProcessor(callback func(letterboxd.Event)) Processor {
	return Processor{
		callback:                          callback,
		watchedDurationByImdbId:           make(map[string]time.Duration),
		playbackStartNotificationByImdbId: make(map[string]PlaybackNotification),
		playbackStopTimeByImdbId:          make(map[string]time.Time),
	}
}

func (p *Processor) ProcessWatchedNotification(notification WatchedNotification) {
	slog.Info(fmt.Sprintf("Processing watched notification %+v", notification))

	var action letterboxd.Action
	if notification.Watched {
		var watchedPercentage = uint(p.watchedDurationByImdbId[notification.ImdbId].Nanoseconds() * 100 / notification.Runtime.Nanoseconds())
		if watchedPercentage >= _MIN_WATCHED_PERCENTAGE {
			action = letterboxd.FilmLogged
		} else {
			action = letterboxd.FilmWatched
		}
	} else {
		action = letterboxd.FilmUnwatched
	}

	delete(p.watchedDurationByImdbId, notification.ImdbId)
	delete(p.playbackStartNotificationByImdbId, notification.ImdbId)
	delete(p.playbackStopTimeByImdbId, notification.ImdbId)

	p.callback(letterboxd.Event{
		ImdbId: notification.ImdbId,
		Action: action,
		Time:   notification.Time,
	})
}

func (p *Processor) ProcessPlaybackNotification(notification PlaybackNotification) {
	slog.Info(fmt.Sprintf("Processing playback notification %+v", notification))

	// TODO: setup DB for permanent storage of partially watched films
	var startNotification, hasStart = p.playbackStartNotificationByImdbId[notification.ImdbId]
	if notification.Playing {
		if !hasStart {
			// Keep earliest playback notification for current session
			p.playbackStartNotificationByImdbId[notification.ImdbId] = notification
		}
		delete(p.playbackStopTimeByImdbId, notification.ImdbId)
	} else {
		if hasStart {
			var watchedDuration = min(
				// Ensure movie was actually watched
				notification.Time.Sub(startNotification.Time),
				// Ensure rewinding/replaying is not included in watched duration
				max(notification.Position - startNotification.Position, 0),
			)
			p.watchedDurationByImdbId[notification.ImdbId] += watchedDuration
			delete(p.playbackStartNotificationByImdbId, notification.ImdbId)
		} else if notification.Time.Sub(p.playbackStopTimeByImdbId[notification.ImdbId]) <= _MAX_DUPLICATE_STOP_PLAYBACK_ELAPSED_TIME {
			slog.Info("Ignoring duplicate playback stop notification")
			return
		} else {
			slog.Warn("Missing playback start time, set total watched duration to current playback position")
			p.watchedDurationByImdbId[notification.ImdbId] = notification.Position
		}
		p.playbackStopTimeByImdbId[notification.ImdbId] = notification.Time

		var positionPercentage = uint(notification.Position.Nanoseconds() * 100 / notification.Runtime.Nanoseconds())
		if positionPercentage >= _MIN_POSITION_PERCENTAGE {
			var watchedPercentage = uint(p.watchedDurationByImdbId[notification.ImdbId].Nanoseconds() * 100 / notification.Runtime.Nanoseconds())
			if watchedPercentage >= _MIN_WATCHED_PERCENTAGE {
				p.callback(letterboxd.Event{
					ImdbId: notification.ImdbId,
					Action: letterboxd.FilmLogged,
					Time:   notification.Time,
				})
			}
			delete(p.watchedDurationByImdbId, notification.ImdbId)
		}
	}
}

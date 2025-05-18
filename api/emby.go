package api

import (
	"time"
	"log/slog"
)

import "github.com/gin-gonic/gin"

import (
	"github.com/computer-geek64/emboxd/notification"
)

const _EMBY_TIME_LAYOUT string = "2006-01-02T15:04:05.0000000Z"

type embyNotification struct {
	Title string `json:"Title"`
	Date  string `json:"Date"`
	Event string `json:"Event"`
	User  struct {
		Name string `json:"Name"`
	} `json:"User"`
	Item struct {
		Type string `json:"Type"`
		RuntimeTicks int64 `json:"RunTimeTicks"`
		ProviderIds  struct {
			Imdb string `json:"Imdb"`
		} `json:"ProviderIds"`
	} `json:"Item"`
	PlaybackInfo struct {
		PlayedToCompletion bool   `json:"PlayedToCompletion"`
		PositionTicks      int64  `json:"PositionTicks"`
		PlaylistIndex      int    `json:"PlaylistIndex"`
		PlaylistLength     int    `json:"PlaylistLength"`
		PlaySessionId      string `json:"PlaySessionId"`
	} `json:"PlaybackInfo"`
}

func convertTicksToDuration(ticks int64) time.Duration {
	return time.Duration(ticks / 10 * int64(time.Microsecond))
}

func (a *Api) postEmbyWebhook(context *gin.Context) {
	var embyNotif embyNotification
	if err := context.BindJSON(&embyNotif); err != nil {
		slog.Error("Malformed webhook notification payload")
		context.AbortWithError(400, err)
		return
	}

	var notificationProcessor, knownEmbyUser = a.notificationProcessorByEmbyUsername[embyNotif.User.Name]
	if !knownEmbyUser {
		// Ignore notifications from unconfigured users
		slog.Debug("No Letterboxd account for Emby user, ignoring notification", slog.Group("emby", "user", embyNotif.User.Name))
		context.AbortWithStatus(200)
		return
	}

	if embyNotif.Item.Type != "Movie" || embyNotif.Item.ProviderIds.Imdb == "" {
		// Only handle movies and valid IMDB entries
		slog.Debug("Media item is not a valid movie, ignoring notification", slog.Group("emby", "user", embyNotif.User.Name, "type", embyNotif.Item.Type), slog.Group("imdb", "id", embyNotif.Item.ProviderIds.Imdb))
		context.AbortWithStatus(200)
		return
	}

	var eventTime, timeErr = time.Parse(_EMBY_TIME_LAYOUT, embyNotif.Date)
	if timeErr != nil {
		slog.Error("Failed to parse time from Emby notification", slog.Group("emby", "user", embyNotif.User.Name, "time", embyNotif.Date))
		context.AbortWithError(400, timeErr)
		return
	}
	var metadata = notification.Metadata{
		Server:   notification.Emby,
		Username: embyNotif.User.Name,
		ImdbId:   embyNotif.Item.ProviderIds.Imdb,
		Time:     eventTime,
	}

	switch embyNotif.Event {
	case "item.markplayed":
		notificationProcessor.ProcessWatchedNotification(notification.WatchedNotification{
			Metadata: metadata,
			Watched:  true,
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "item.markunplayed":
		notificationProcessor.ProcessWatchedNotification(notification.WatchedNotification{
			Metadata: metadata,
			Watched:  false,
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "playback.start", "playback.unpause":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Playing:  true,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "playback.stop", "playback.pause":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Playing:  false,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})

		if embyNotif.PlaybackInfo.PlayedToCompletion {
			notificationProcessor.ProcessWatchedNotification(notification.WatchedNotification{
				Metadata: metadata,
				Watched:  true,
				Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
			})
		}
	default:
		context.AbortWithStatus(400)
		return
	}
}

func (a *Api) setupEmbyRoutes() {
	var embyRouter = a.router.Group("/emby")
	embyRouter.POST("/webhook", a.postEmbyWebhook)
}

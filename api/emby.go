package api

import (
	"time"
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
		RuntimeTicks int64 `json:"PositionTicks"`
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
		context.AbortWithError(400, err)
		return
	}

	var notificationProcessor, knownEmbyUser = a.notificationProcessorByEmbyUsername[embyNotif.User.Name]
	if !knownEmbyUser {
		// Ignore notifications from unconfigured users
		context.AbortWithStatus(200)
		return
	}

	if embyNotif.Item.ProviderIds.Imdb == "" {
		// Only handle valid IMDB entries
		context.AbortWithStatus(200)
		return
	}

	var eventTime, timeErr = time.Parse(_EMBY_TIME_LAYOUT, embyNotif.Date)
	if timeErr != nil {
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
		})
	case "item.markunplayed":
		notificationProcessor.ProcessWatchedNotification(notification.WatchedNotification{
			Metadata: metadata,
			Watched:  false,
		})
	case "playback.start":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Action:   notification.Start,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "playback.unpause":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Action:   notification.Start,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "playback.pause":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Action:   notification.Start,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})
	case "playback.stop":
		notificationProcessor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Action:   notification.Start,
			Position: convertTicksToDuration(embyNotif.PlaybackInfo.PositionTicks),
			Runtime:  convertTicksToDuration(embyNotif.Item.RuntimeTicks),
		})

		if embyNotif.PlaybackInfo.PlayedToCompletion {
			notificationProcessor.ProcessWatchedNotification(notification.WatchedNotification{
				Metadata: metadata,
				Watched:  true,
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

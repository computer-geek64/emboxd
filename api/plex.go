package api

import (
	"log/slog"
	"time"

	"github.com/computer-geek64/emboxd/notification"
	"github.com/gin-gonic/gin"
)

// Plex webhook payload structure (simplified for movie events)
type plexNotification struct {
	Event   string `json:"event"`
	Account struct {
		Title string `json:"title"`
	} `json:"Account"`
	Metadata struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Guid     string `json:"guid"`
		Duration int64  `json:"duration"`
	} `json:"Metadata"`
	Player struct {
		State string `json:"state"`
	} `json:"Player"`
	Server struct {
		Title string `json:"title"`
	} `json:"Server"`
	EventTime int64 `json:"eventTime"`
}

func parsePlexImdbId(guid string) string {
	// Plex IMDB GUIDs are in the form: "imdb://tt1234567"
	const prefix = "imdb://"
	if len(guid) > len(prefix) && guid[:len(prefix)] == prefix {
		return guid[len(prefix):]
	}
	return ""
}

func (a *Api) postPlexWebhook(context *gin.Context) {
	var plexNotif plexNotification
	if err := context.BindJSON(&plexNotif); err != nil {
		slog.Error("Malformed Plex webhook notification payload")
		context.AbortWithError(400, err)
		return
	}

	// Only handle movies with IMDB id
	if plexNotif.Metadata.Type != "movie" {
		context.AbortWithStatus(200)
		return
	}
	imdbId := parsePlexImdbId(plexNotif.Metadata.Guid)
	if imdbId == "" {
		context.AbortWithStatus(200)
		return
	}

	username := plexNotif.Account.Title
	processor, ok := a.notificationProcessorByPlexUsername[username]
	if !ok {
		slog.Debug("No Letterboxd account for Plex user, ignoring notification", slog.Group("plex", "user", username))
		context.AbortWithStatus(200)
		return
	}
	eventTime := time.Unix(plexNotif.EventTime, 0)
	metadata := notification.Metadata{
		Server:   notification.Plex,
		Username: username,
		ImdbId:   imdbId,
		Time:     eventTime,
	}

	switch plexNotif.Event {
	case "media.scrobble":
		processor.ProcessWatchedNotification(notification.WatchedNotification{
			Metadata: metadata,
			Watched:  true,
			Runtime:  time.Duration(plexNotif.Metadata.Duration) * time.Millisecond,
		})
	case "media.play", "media.resume":
		processor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Playing:  true,
			Position: 0, // Plex does not provide position in webhook
			Runtime:  time.Duration(plexNotif.Metadata.Duration) * time.Millisecond,
		})
	case "media.pause", "media.stop":
		processor.ProcessPlaybackNotification(notification.PlaybackNotification{
			Metadata: metadata,
			Playing:  false,
			Position: 0, // Plex does not provide position in webhook
			Runtime:  time.Duration(plexNotif.Metadata.Duration) * time.Millisecond,
		})
	default:
		context.AbortWithStatus(400)
		return
	}

	context.Status(200)
}

func (a *Api) setupPlexRoutes() {
	plexRouter := a.router.Group("/plex")
	plexRouter.POST("/webhook", a.postPlexWebhook)
}

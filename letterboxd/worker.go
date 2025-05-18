package letterboxd

import (
	"fmt"
	"log/slog"
	"time"
)

const _EVENT_BUFFER_SIZE int = 10

type Action int

const (
	FilmUnwatched Action = iota
	FilmWatched
	FilmLogged
)

type Event struct {
	ImdbId string
	Action Action
	Time   time.Time
}

type Worker struct {
	debouncer
	user    User
	channel chan Event
}

func NewWorker(username string, password string) Worker {
	var channel = make(chan Event, _EVENT_BUFFER_SIZE)
	return Worker{
		debouncer: newDebouncer(
			channel,
		),
		user: NewUser(
			username,
			password,
		),
		channel: channel,
	}
}

func (w *Worker) HandleEvent(event Event) {
	w.debounce(event)
}

func (w *Worker) Start() {
	go w.run()
}

func (w *Worker) run() {
	w.user.Login()

	for {
		var event = <-w.channel

		switch event.Action {
		case FilmWatched, FilmUnwatched:
			w.user.SetFilmWatched(event.ImdbId, event.Action == FilmWatched)
		case FilmLogged:
			w.user.LogFilmWatched(event.ImdbId)
		default:
			panic(fmt.Sprintf("Unknown event action %d", event.Action))
		}

		slog.Info(fmt.Sprintf("Finished processing event %+v", event))
	}
}

package letterboxd

import (
	"container/list"
	"sync"
	"time"
)

const _QUIET_PERIOD = 30 * time.Second

type debouncer struct {
	queue                    *list.List
	removableElementByImdbId map[string]*list.Element
	loggedImdbIds            map[string]bool
	lock                     sync.Mutex
	channel                  chan Event
}

func newDebouncer(channel chan Event) debouncer {
	return debouncer{
		queue:                    list.New(),
		removableElementByImdbId: make(map[string]*list.Element),
		loggedImdbIds:            make(map[string]bool),
		channel:                  channel,
	}
}

func (d *debouncer) debounce(event Event) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if element, ok := d.removableElementByImdbId[event.ImdbId]; ok {
		d.queue.Remove(element)
	}

	if event.Action == FilmLogged {
		d.loggedImdbIds[event.ImdbId] = true
	}

	// TODO: fix
	d.channel <- event
}

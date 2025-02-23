package letterboxd

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
)

func (u User) SetFilmWatched(imdbId string, watched bool) {
	var url = fmt.Sprintf("https://letterboxd.com/imdb/%s", imdbId)
	var page = u.newPage(url)
	defer page.Close()

	// Reauthenticate if necessary
	if !u.isLoggedIn(page) {
		slog.Warn("Not logged in, authenticating...")

		u.Login()
		if _, err := page.Reload(); err != nil {
			slog.Warn(fmt.Sprintf("Page %s took too long to load", url))
		}
	}

	// Allow watched information to populate
	time.Sleep(3 * time.Second)

	var watchedLocator = page.Locator("span.action-large.-watch .action.-watch")
	var classes, watchedLocatorErr = watchedLocator.GetAttribute("class")
	if watchedLocatorErr != nil {
		panic(watchedLocatorErr)
	}

	if slices.Contains(strings.Split(classes, " "), "-on") == watched {
		// Film already marked with desired watch state
		slog.Info(fmt.Sprintf("Film %s is already marked as watched = %t", imdbId, watched))
	} else {
		// Toggle film watched status
		if err := watchedLocator.Click(); err != nil {
			panic(err)
		}
		time.Sleep(3 * time.Second)
	}
}

func (u User) LogFilmWatched(imdbId string, date ...time.Time) {
	if len(date) == 0 {
		date = append(date, time.Now())
	}

	var url = fmt.Sprintf("https://letterboxd.com/imdb/%s", imdbId)
	var page = u.newPage(url)
	defer page.Close()

	// Reauthenticate if necessary
	if !u.isLoggedIn(page) {
		slog.Warn("Not logged in, authenticating...")

		u.Login()
		if _, err := page.Reload(); err != nil {
			slog.Warn(fmt.Sprintf("Page %s took too long to load", url))
		}
	}

	if err := page.Locator("button.add-this-film").Click(); err != nil {
		panic(err)
	}

	// Wait for form to load
	var saveLocator = page.Locator("div#diary-entry-form-modal button.button.-action.button-action")
	if err := saveLocator.WaitFor(); err != nil {
		panic(err)
	}

	// Fill form and save log entry
	var javascriptSetDate = fmt.Sprintf("document.querySelector('input#frm-viewing-date-string').value = '%s'", date[0].Format(time.DateOnly))
	if _, err := page.Evaluate(javascriptSetDate, nil); err != nil {
		panic(err)
	}
	if err := saveLocator.Click(); err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
}

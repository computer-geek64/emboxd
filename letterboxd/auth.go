package letterboxd

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

import "github.com/playwright-community/playwright-go"

func (u User) isLoggedIn(page ...playwright.Page) bool {
	if len(page) == 0 {
		// Create new page
		page = append(page, u.newPage("https://letterboxd.com"))
		defer page[0].Close()
	}

	classes, err := page[0].Locator("body").GetAttribute("class")
	if err != nil {
		panic(err)
	}

	return slices.Contains(strings.Split(classes, " "), "logged-in")
}

func (u User) Login() {
	page := u.newPage("https://letterboxd.com/sign-in/")
	defer page.Close()

	if page.URL() == "https://letterboxd.com" {
		slog.Warn("Already logged in")
		return
	}

	// Fill out login form
	if err := page.Locator("input#field-username").Fill(u.username); err != nil {
		panic(err)
	}
	if err := page.Locator("input#field-password").Fill(u.password); err != nil {
		panic(err)
	}

	// Remember me (optional, ignore errors)
	_ = page.Locator("input.js-remember").Check()

	// SUBMIT VIA ENTER (robust, avoids overlays / cloudflare issues)
	if err := page.Locator("input#field-password").Press("Enter"); err != nil {
		panic(err)
	}

	// Wait for navigation / login completion
	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})

	// Confirm logged in
	if err := page.Locator("body.logged-in").WaitFor(); err == nil {
		slog.Info(fmt.Sprintf("Logged in as %s", u.username))
	} else {
		panic(err)
	}
}
	// Wait for logged in status
	if err := page.Locator("body.logged-in").WaitFor(); err == nil {
		slog.Info(fmt.Sprintf("Logged in as %s", u.username))
	} else {
		panic(err)
	}
}

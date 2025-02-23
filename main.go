package main

import "os"

import (
	"github.com/computer-geek64/emboxd/api"
	"github.com/computer-geek64/emboxd/config"
	"github.com/computer-geek64/emboxd/letterboxd"
	"github.com/computer-geek64/emboxd/notification"
)

func main() {
	var configFilename = "config/config.yaml"
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-c" || os.Args[i] == "--config" {
			i++
			configFilename = os.Args[i]
			break
		}
	}
	var conf = config.Load(configFilename)

	var notificationProcessorByEmbyUsername = make(map[string]*notification.Processor, len(conf.Users))
	for _, user := range conf.Users {
		var letterboxdWorker = letterboxd.NewWorker(user.Letterboxd.Username, user.Letterboxd.Password)
		letterboxdWorker.Start()

		var notificationProcessor = notification.NewProcessor(letterboxdWorker.HandleEvent)

		notificationProcessorByEmbyUsername[user.Emby.Username] = &notificationProcessor
	}

	var app = api.New(notificationProcessorByEmbyUsername)
	app.Run(80)
}

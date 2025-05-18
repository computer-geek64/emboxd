package main

import (
	"flag"
)

import (
	"github.com/computer-geek64/emboxd/api"
	"github.com/computer-geek64/emboxd/config"
	"github.com/computer-geek64/emboxd/letterboxd"
	"github.com/computer-geek64/emboxd/logging"
	"github.com/computer-geek64/emboxd/notification"
)

func main() {
	var verbose bool
	var configFilename string
	flag.BoolVar(&verbose, "v", false, "Enable debug logging")
	flag.BoolVar(&verbose, "verbose", false, "Enable debug logging")
	flag.StringVar(&configFilename, "c", "config/config.yaml", "Path to configuration file")
	flag.StringVar(&configFilename, "config", "config/config.yaml", "Path to configuration file")
	flag.Parse()

	logging.Configure(verbose)
	var conf = config.Load(configFilename)

	var notificationProcessorByEmbyUsername = make(map[string]*notification.Processor, len(conf.Users))
	var letterboxdWorkers = make(map[string]*letterboxd.Worker, len(conf.Users))
	for _, user := range conf.Users {
		var letterboxdWorker, workerExists = letterboxdWorkers[user.Letterboxd.Username]
		if !workerExists {
			var worker = letterboxd.NewWorker(user.Letterboxd.Username, user.Letterboxd.Password)
			worker.Start()
			letterboxdWorker = &worker
			letterboxdWorkers[user.Letterboxd.Username] = letterboxdWorker
		}

		var notificationProcessor = notification.NewProcessor(letterboxdWorker.HandleEvent)
		notificationProcessorByEmbyUsername[user.Emby.Username] = &notificationProcessor
	}

	var app = api.New(notificationProcessorByEmbyUsername)
	app.Run(80)
}

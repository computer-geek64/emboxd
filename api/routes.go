package api

import (
	"fmt"
	"log/slog"
)

import "github.com/gin-gonic/gin"

import "github.com/computer-geek64/emboxd/notification"

type Api struct {
	router                              *gin.Engine
	notificationProcessorByEmbyUsername map[string]*notification.Processor
}

func New(notificationProcessorByEmbyUsername map[string]*notification.Processor) Api {
	gin.SetMode(gin.ReleaseMode)
	return Api{
		router:                              gin.Default(),
		notificationProcessorByEmbyUsername: notificationProcessorByEmbyUsername,
	}
}

func (a *Api) getRoot(context *gin.Context) {
	context.String(200, "Welcome to EmBoxd!")
}

func (a *Api) setupRoutes() {
	a.setupEmbyRoutes()

	a.router.GET("/", a.getRoot)
}

func (a *Api) Run(port int) {
	a.setupRoutes()

	slog.Info("Starting Gin Server")
	a.router.Run(fmt.Sprintf(":%d", port))
}

package app

import (
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/endpoint"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
)

type App struct {
	e  *endpoint.Endpoint
	ry *repository.Repository
	s  *service.Service
	rr *gin.Engine
	c  *config.Config
}

func New() (*App, error) {
	a := &App{}

	a.c = config.New()

	a.ry = repository.New()

	a.s = service.New(a.ry)

	a.e = endpoint.New(a.s)

	a.rr = gin.Default()

	a.rr.POST("/", a.e.CreateShortened)
	a.rr.POST("/api/shorten", a.e.JsonCreateShortened)
	a.rr.GET("/:id", a.e.GetOriginalURL)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running. PORT :" + a.c.ServerAddress)
	a.rr.Run(":" + a.c.ServerAddress)

	return nil
}

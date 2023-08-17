package app

import (
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/endpoint"
	"github.com/april1858/shortener-gin/internal/app/middleware"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
)

type App struct {
	e  *endpoint.Endpoint
	rp *repository.Repository
	s  *service.Service
	r  *gin.Engine
	c  *config.Config
	mw *middleware.MW
}

func New() (*App, error) {
	a := &App{}

	a.c = config.New()

	a.rp = repository.New(a.c)

	a.s = service.New(a.rp)

	a.e = endpoint.New(a.s)

	a.mw = middleware.New()

	a.r = gin.Default()

	a.r.Use(a.mw.GZIP())
	a.r.POST("/", a.e.CreateShortened)
	a.r.POST("/api/shorten", a.e.JSONCreateShortened)
	a.r.GET("/:id", a.e.GetOriginalURL)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")
	a.r.Run(a.c.ServerAddress)

	return nil
}

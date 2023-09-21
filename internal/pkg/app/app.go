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
	e      *endpoint.Endpoint
	repo   *repository.Repository
	s      *service.Service
	route  *gin.Engine
	config *config.Config
	mw     *middleware.MW
}

func New() (*App, error) {
	a := &App{}

	a.config = config.New()

	a.repo = repository.New(*a.config)

	a.s = service.New(a.repo)

	a.e = endpoint.New(a.s)

	a.mw = middleware.New()

	gin.SetMode(gin.ReleaseMode)
	a.route = gin.Default()

	a.route.Use(a.mw.Cookie(), a.mw.GZIP())
	a.route.POST("/", a.e.CreateShortened)
	a.route.POST("/api/shorten", a.e.JSONCreateShortened)
	a.route.POST("/api/shorten/batch", a.e.CreateShortenedBatch)
	a.route.GET("/:id", a.e.GetOriginalURL)
	a.route.GET("/api/user/urls", a.e.GetAllUID)
	a.route.GET("/ping", a.e.Ping)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")
	a.route.Run(a.config.ServerAddress)

	return nil
}

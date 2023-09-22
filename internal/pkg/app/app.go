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
	endpoint *endpoint.Endpoint
	repo     *repository.Repository
	service  *service.Service
	route    *gin.Engine
	config   *config.Config
	mw       *middleware.MW
}

func New() (*App, error) {
	a := &App{}

	a.config = config.New()

	a.repo = repository.New(*a.config)

	a.service = service.New(a.repo)

	a.endpoint = endpoint.New(a.service)

	a.mw = middleware.New()

	gin.SetMode(gin.ReleaseMode)
	a.route = gin.Default()

	a.route.Use(a.mw.Cookie(), a.mw.GZIP())
	a.route.POST("/", a.endpoint.CreateShortened)
	a.route.POST("/api/shorten", a.endpoint.JSONCreateShortened)
	a.route.POST("/api/shorten/batch", a.endpoint.CreateShortenedBatch)
	a.route.GET("/:id", a.endpoint.GetOriginalURL)
	a.route.GET("/api/user/urls", a.endpoint.GetAllUID)
	a.route.GET("/ping", a.endpoint.Ping)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")
	a.route.Run(a.config.ServerAddress)

	return nil
}

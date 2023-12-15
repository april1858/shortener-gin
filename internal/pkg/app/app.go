package app

import (
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/endpoint"
	"github.com/april1858/shortener-gin/internal/app/entity"
	"github.com/april1858/shortener-gin/internal/app/middleware"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type App struct {
	endpoint   *endpoint.Endpoint
	repository repository.Repository
	service    *service.Service
	route      *gin.Engine
	config     *config.Config
	mw         *middleware.MW
}

func New() (*App, error) {
	var err error
	var ch chan entity.ChData
	a := &App{}

	a.config = config.New()

	a.repository, ch, err = repository.New(a.config)
	if err != nil {
		fmt.Println("error from repository", err)
	}

	a.service, ch = service.New(a.repository, ch)

	a.endpoint = endpoint.New(a.service, ch)

	a.mw = middleware.New()

	gin.SetMode(gin.ReleaseMode)
	a.route = gin.Default()
	pprof.Register(a.route)

	a.route.Use(a.mw.Cookie(), a.mw.GZIP())
	a.route.POST("/", a.endpoint.CreateShortened)
	a.route.POST("/api/shorten", a.endpoint.JSONCreateShortened)
	a.route.POST("/api/shorten/batch", a.endpoint.CreateShortenedBatch)
	a.route.GET("/:id", a.endpoint.GetOriginalURL)
	a.route.GET("/api/user/urls", a.endpoint.GetAllUID)
	a.route.GET("/ping", a.endpoint.Ping)
	a.route.DELETE("/api/user/urls", a.endpoint.Delete)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")
	a.route.Run(a.config.ServerAddress)

	return nil
}

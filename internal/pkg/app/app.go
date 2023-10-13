package app

import (
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/endpoint"
	"github.com/april1858/shortener-gin/internal/app/middleware"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
)

type App struct {
	endpoint *endpoint.Endpoint
	//repoDB     *repository.DB
	//repoFile   *repository.File
	//repoMemory *repository.Memory
	service *service.Service
	route   *gin.Engine
	config  *config.Config
	mw      *middleware.MW
}

func New() (*App, error) {
	//var err error
	a := &App{}

	a.config = config.New()
	/*
		switch {
		case a.config.DatabaseDsn != "":
			a.repoDB, err = repository.NewDBStorage(a.config.DatabaseDsn)
			if err != nil {
				return nil, err
			}
			a.service = service.New(a.repoDB)
		case a.config.FileStoragePath != "":
			a.repoFile = repository.NewFileStorage(a.config.FileStoragePath)
			a.service = service.New(a.repoFile)
		default:
			a.repoMemory = repository.NewMemStorage()
			a.service = service.New(a.repoMemory)
		}
	*/
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

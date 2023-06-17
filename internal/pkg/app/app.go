package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/endpoint"
	"github.com/april1858/shortener-gin/internal/app/service"
)

type App struct {
	e *endpoint.Endpoint
	rep *repository.Repository
	s *service.Service
	r *gin.Engine
}

func New() (*App, error) {
	a := &App{}

	a.rep = repository.New()

	a.s = service.New(a.rep)

	a.e = endpoint.New(a.s)

	a.r = gin.Default()

	a.r.POST("/", a.e.CreateShortened)
	a.r.GET("/:id", a.e.GetOriginalURL)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")
	a.r.Run(":8080")

	return nil
}

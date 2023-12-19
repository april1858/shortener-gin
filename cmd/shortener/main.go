package main

import (
	"log"

	"github.com/april1858/shortener-gin/internal/pkg/app"
)

// @Title Shortener-Gin API
// @Description Сервис сокращения URL.
// @Version 1.0

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}

}

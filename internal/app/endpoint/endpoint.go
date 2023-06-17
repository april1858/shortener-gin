package endpoint

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var BaseURL string = "http://localhost:8080/"

var data = make(map[string]string)

type Service interface {
	CreatorShortened(string) string
	FindOriginalURL(string) (string, error)
}

type Endpoint struct {
	S Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		S: s,
	}
}

func (e *Endpoint) CreateShortened(c *gin.Context) {
	originalURL, _ := c.GetRawData()
	_, err := url.ParseRequestURI(string(originalURL))
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain", []byte("Не правильный URL"))
	} else {
		shortened := e.S.CreatorShortened(string(originalURL))
		c.Data(http.StatusCreated, "text/plain", []byte(BaseURL+shortened))
	}
}

func (e *Endpoint) GetOriginalURL(c *gin.Context) {
	shortened := c.Param("id")
	answer, err := e.S.FindOriginalURL(shortened)
	if err != nil {
		s := fmt.Sprintf("Ошибка - %v", err)
		c.Data(http.StatusBadRequest, "text/plain", []byte(s))
	} else {
		c.Redirect(http.StatusTemporaryRedirect, answer)
	}
}

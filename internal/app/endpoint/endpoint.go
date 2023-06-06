package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var BaseURL string = "http://localhost:8080/"

var data = make(map[string]string)

type Service interface {
	CreatorShortened() string
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}



func (e *Endpoint) CreateShortened(c *gin.Context) {
	body, _ := c.GetRawData()
	shortened := e.s.CreatorShortened()
	data[shortened] = string(body)
	c.Data(http.StatusCreated, "", []byte(BaseURL + shortened))
}

func (e *Endpoint) GetOriginalURL(c *gin.Context) {
	shortened := c.Param("id")
	answer := data[shortened]
	c.Redirect(http.StatusTemporaryRedirect, answer)
}

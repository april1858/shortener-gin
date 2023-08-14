package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/april1858/shortener-gin/internal/app/config"

	"github.com/gin-gonic/gin"
)

//var baseURL = config.BaseURL

//var data = make(map[string]string)

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
		c.Data(http.StatusCreated, "text/plain", []byte(config.Cnf.BaseURL+shortened))
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

func (e *Endpoint) JsonCreateShortened(c *gin.Context) {
	var shortened string
	objQuery := make(map[string]string)
	requestBody, _ := c.GetRawData()

	if err := json.Unmarshal(requestBody, &objQuery); err != nil {
		return
	}

	_, err := url.ParseRequestURI(objQuery["url"])
	if err != nil {
		c.Data(http.StatusBadRequest, "application/json", []byte("Не правильный URL"))
	} else {
		shortened = e.S.CreatorShortened(objQuery["url"])
	}

	answerStruct := map[string]string{"result": config.Cnf.BaseURL + shortened}
	answer, err := json.Marshal(answerStruct)
	if err != nil {
		return
	}

	c.Data(http.StatusCreated, "application/json", answer)
}

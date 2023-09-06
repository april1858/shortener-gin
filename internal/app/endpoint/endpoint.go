package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/april1858/shortener-gin/internal/app/config"

	"github.com/gin-gonic/gin"
)

type Redirect struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type Service interface {
	CreatorShortened(string, string) (string, error)
	FindOriginalURL(string) (string, error)
	FindByUID(uid string) ([]string, error)
	Ping() (string, error)
	CreatorShortenedBatch([]map[string]string, string) []string
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
	contentType := c.GetHeader("Accept")
	var status int = http.StatusCreated
	originalURL, _ := c.GetRawData()
	_, err := url.ParseRequestURI(string(originalURL))
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain", []byte("Не правильный URL"))
	} else {
		shortened, err := e.S.CreatorShortened(string(originalURL), c.MustGet("UID").(string))
		if err != nil {
			status = http.StatusConflict
		}
		c.Data(status, contentType, []byte(config.Cnf.BaseURL+shortened))
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

func (e *Endpoint) GetAllUID(c *gin.Context) {
	fmt.Println("c.MustGet().(string)", c.MustGet("UID").(string))
	sliceAll, err := e.S.FindByUID(c.MustGet("UID").(string))
	if err != nil {
		s := fmt.Sprintf("Ошибка - %v", err)
		c.Data(http.StatusNoContent, "text/plain application/json", []byte(s))
	} else {
		var redirect = make([]Redirect, 0, 1)
		var r Redirect
		for _, value := range sliceAll {
			var v = strings.Fields(value)
			r.ShortURL = config.Cnf.BaseURL + v[0]
			r.OriginalURL = v[1]
			redirect = append(redirect, r)
		}
		answer, err := json.Marshal(redirect)
		if err != nil {
			return
		}
		c.Header("WWW-Authenticate", `Basic realm="api"`)
		c.Data(http.StatusOK, "text/plain application/json", answer)
	}
}

func (e *Endpoint) JSONCreateShortened(c *gin.Context) {
	var shortened string
	var status int = http.StatusCreated
	objQuery := make(map[string]string)
	requestBody, _ := c.GetRawData()

	if err := json.Unmarshal(requestBody, &objQuery); err != nil {
		return
	}

	_, err := url.ParseRequestURI(objQuery["url"])
	if err != nil {
		c.Data(http.StatusBadRequest, "application/json", []byte("Не правильный URL"))
	} else {
		shortened, err = e.S.CreatorShortened(objQuery["url"], c.MustGet("UID").(string))
		if err != nil {
			status = http.StatusConflict
		}
	}

	answerStruct := map[string]string{"result": config.Cnf.BaseURL + shortened}
	answer, err := json.Marshal(answerStruct)
	if err != nil {
		return
	}

	c.Data(status, "application/json", answer)
}

func (e *Endpoint) Ping(c *gin.Context) {
	_, err := e.S.Ping()
	if err != nil {
		c.Data(http.StatusInternalServerError, "", nil)
	}
	c.Data(http.StatusOK, "", nil)
}

func (e *Endpoint) CreateShortenedBatch(c *gin.Context) {
	objQuery := make([]map[string]string, 0)
	requestBody, _ := c.GetRawData()

	if err := json.Unmarshal(requestBody, &objQuery); err != nil {
		fmt.Println("err - ", err)
		return
	}
	answer := e.S.CreatorShortenedBatch(objQuery, c.MustGet("UID").(string))
	for i, v := range objQuery {
		delete(v, "original_url")
		v["short_url"] = config.Cnf.BaseURL + strings.Fields(answer[i])[0]
	}
	answer1, err := json.Marshal(objQuery)

	if err != nil {
		return
	}
	c.Data(http.StatusCreated, "application/json", []byte(answer1))

}

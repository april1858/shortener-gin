package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Redirect struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Service interface {
	CreatorShortened(*gin.Context, string) (string, error)
	FindOriginalURL(*gin.Context, string) (string, error)
	FindByUID(*gin.Context) ([]string, error)
	Ping() (string, error)
	CreatorShortenedBatch(*gin.Context, []map[string]string) ([]string, error)
	Delete(*gin.Context, chan repository.S)
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

func (e *Endpoint) CreateShortened(ctx *gin.Context) {
	contentType := "text/plain"
	var status int = http.StatusCreated
	originalURL, _ := ctx.GetRawData()
	_, err := url.ParseRequestURI(string(originalURL))
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte("Не правильный URL"))
	} else {
		shortened, err := e.s.CreatorShortened(ctx, string(originalURL))
		if err != nil {
			status = http.StatusConflict
		}
		ctx.Data(status, contentType, []byte(config.BURL+shortened))
	}
}

func (e *Endpoint) GetOriginalURL(ctx *gin.Context) {
	shortened := ctx.Param("id")
	answer, err := e.s.FindOriginalURL(ctx, shortened)
	if answer == "" {
		ctx.Data(http.StatusBadRequest, "text/plain", []byte("Not found"))
	}
	if answer == "deleted" {
		ctx.Data(http.StatusGone, "text/plain", []byte(answer))
	}
	if err != nil {
		s := fmt.Sprintf("Ошибка - %v", err)
		ctx.Data(http.StatusBadRequest, "text/plain", []byte(s))
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, answer)
	}
}

func (e *Endpoint) GetAllUID(ctx *gin.Context) {
	sliceAll, err := e.s.FindByUID(ctx)
	if err != nil {
		s := fmt.Sprintf("Ошибка - %v", err)
		ctx.Data(http.StatusNoContent, "text/plain application/json", []byte(s))
	} else {
		var redirect = make([]Redirect, 0, 1)
		var r Redirect
		for _, value := range sliceAll {
			var v = strings.Fields(value)
			r.ShortURL = config.BURL + v[0]
			r.OriginalURL = v[1]
			redirect = append(redirect, r)
		}
		answer, err := json.Marshal(redirect)
		if err != nil {
			return
		}
		ctx.Header("WWW-Authenticate", `Basic realm="api"`)
		ctx.Data(http.StatusOK, "text/plain application/json", answer)
	}
}

func (e *Endpoint) JSONCreateShortened(ctx *gin.Context) {
	var shortened string
	var status int = http.StatusCreated
	objQuery := make(map[string]string)
	requestBody, _ := ctx.GetRawData()

	if err := json.Unmarshal(requestBody, &objQuery); err != nil {
		return
	}

	_, err := url.ParseRequestURI(objQuery["url"])
	if err != nil {
		ctx.Data(http.StatusBadRequest, "application/json", []byte("Не правильный URL"))
	} else {
		shortened, err = e.s.CreatorShortened(ctx, objQuery["url"])
		if err != nil {
			status = http.StatusConflict
		}
	}

	answerStruct := map[string]string{"result": config.BURL + shortened}
	answer, err := json.Marshal(answerStruct)
	if err != nil {
		return
	}

	ctx.Data(status, "application/json", answer)
}

func (e *Endpoint) Ping(ctx *gin.Context) {
	message, err := e.s.Ping()
	if err != nil {
		ctx.Data(http.StatusInternalServerError, "", nil)
	}
	ctx.Data(http.StatusOK, "", []byte(message))
}

func (e *Endpoint) CreateShortenedBatch(ctx *gin.Context) {
	objQuery := make([]map[string]string, 0)
	requestBody, _ := ctx.GetRawData()

	if err := json.Unmarshal(requestBody, &objQuery); err != nil {
		ctx.Data(http.StatusCreated, "application/json", []byte(err.Error()))
	}
	answer, err := e.s.CreatorShortenedBatch(ctx, objQuery)
	if err != nil {
		ctx.Data(http.StatusCreated, "application/json", []byte(err.Error()))
	}
	for i, v := range objQuery {
		delete(v, "original_url")
		v["short_url"] = config.BURL + strings.Fields(answer[i])[0]
	}
	answer1, err := json.Marshal(objQuery)

	if err != nil {
		ctx.Data(http.StatusCreated, "application/json", []byte(err.Error()))
	}
	ctx.Data(http.StatusCreated, "application/json", []byte(answer1))
}

func (e *Endpoint) Delete(ctx *gin.Context) {
	uid := ctx.MustGet("UID").(string)
	c := make(chan repository.S)
	remove := make([]string, 1)
	requestBody, _ := ctx.GetRawData()

	if err := json.Unmarshal(requestBody, &remove); err != nil {
		ctx.Data(http.StatusCreated, "application/json", []byte(err.Error()))
	}
	s := repository.S{UID: uid, Data: remove}
	fmt.Println("remove - ", remove)
	go func(cc chan repository.S) {
		cc <- s
	}(c)
	e.s.Delete(ctx, c)
	ctx.Data(http.StatusAccepted, "application/json", []byte("OK"))

}

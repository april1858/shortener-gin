package endpoint

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestCreateSortened(t *testing.T) {
	rep := repository.New()
	s := service.New(rep)
	ep := New(s)
	r := SetUpRouter()
	r.POST("/", ep.CreateShortened)

	type want struct {
		code        int
		originalURL string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "first",
			want: want{
				code:        201,
				originalURL: "http://s-s.ru/123qweewqwe1313werfsw43we/ertfdsgsdfggfsdfgsdfgsdfgdgsdg",
			},
		},
		{
			name: "second",
			want: want{
				code:        400,
				originalURL: "oh no",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := []byte(tt.want.originalURL)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(ms))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestGetOriginalURL(t *testing.T) {

	rep := repository.New()
	(*rep).M["1234567"] = "http://s-s.ru"
	s := service.New(rep)
	ep := New(s)
	r := SetUpRouter()
	r.GET("/:id", ep.GetOriginalURL)

	type want struct {
		code         int
		shortenedURL string
		originalURL  string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "f1",
			want: want{
				code:         307,
				shortenedURL: "1234567",
				originalURL:  "http://s-s.ru",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/"+tt.want.shortenedURL, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			fmt.Println(w.Header().Get("Location"))
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.originalURL, w.Header().Get("Location"))
		})
	}
}

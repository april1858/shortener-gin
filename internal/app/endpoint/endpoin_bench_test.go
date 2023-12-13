package endpoint

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/april1858/shortener-gin/internal/app/middleware"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouterBench() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	rep := repository.NewMemStorage()
	addMemory := []repository.ES{
		{Short: "1234567", Original: "http://s-s1.ru", UID: "1", Condition: true},
		{Short: "2345678", Original: "http://s-s2.ru", UID: "1", Condition: true},
		{Short: "3456789", Original: "http://s-s3.ru./articles/go/testirovanie-http-hendlerov-v-go/", UID: "1", Condition: true},
	}
	rep.Memory = append(rep.Memory, addMemory...)
	service, ch := service.New(rep, ch)
	endpoint := New(service, ch)
	mw := middleware.New()

	router.Use(mw.Cookie(), mw.GZIP())
	router.POST("/", endpoint.CreateShortened)
	router.GET("/:id", endpoint.GetOriginalURL)

	return router
}

func BenchmarkCreateSortened(b *testing.B) {
	b.StopTimer()
	r := SetUpRouter()
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
				originalURL: "http://asdfghh.ru/123qweewqwe1313werfsw43we/ertfdsgsdfggfsdfgsdfgsdfgdgsdg",
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
	b.StartTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {

			ms := []byte(tt.want.originalURL)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(ms))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(b, tt.want.code, w.Code)
		})
	}
}

func BenchmarkGetOriginalURL(b *testing.B) {
	b.StopTimer()
	r := SetUpRouter()

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
			name: "test1",
			want: want{
				code:         307,
				shortenedURL: "1234567",
				originalURL:  "http://s-s1.ru",
			},
		},
		{
			name: "test2",
			want: want{
				code:         307,
				shortenedURL: "2345678",
				originalURL:  "http://s-s2.ru",
			},
		},
		{
			name: "test3",
			want: want{
				code:         307,
				shortenedURL: "3456789",
				originalURL:  "http://s-s3.ru./articles/go/testirovanie-http-hendlerov-v-go/",
			},
		},
		{
			name: "test4",
			want: want{
				code:         404,
				shortenedURL: "",
				originalURL:  "",
			},
		},
		{
			name: "test5",
			want: want{
				code:         400,
				shortenedURL: "1233456789",
				originalURL:  "",
			},
		},
	}
	b.StartTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			req, _ := http.NewRequest("GET", "/"+tt.want.shortenedURL, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(b, tt.want.code, w.Code)
			assert.Equal(b, tt.want.originalURL, w.Header().Get("Location"))
		})
	}
}

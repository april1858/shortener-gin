package endpoint

/*
import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/april1858/shortener-gin/internal/app/middleware"
	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/april1858/shortener-gin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	rep := repository.NewMemStorage()
	rep.Memory = append(rep.Memory, repository.ES{Short: "1234567", Original: "http://s-s.ru", UID: "1", Condition: true})
	service, ch := service.New(rep, ch)
	endpoint := New(service, ch)
	mw := middleware.New()

	router.Use(mw.Cookie(), mw.GZIP())
	router.POST("/", endpoint.CreateShortened)
	router.GET("/:id", endpoint.GetOriginalURL)

	return router
}

func TestCreateSortened(t *testing.T) {
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
*/

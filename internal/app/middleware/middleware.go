package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

type MW struct {
}

func New() *MW {
	return &MW{}
}

func (mv *MW) GZIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			io.WriteString(c.Writer, err.Error())
			return
		}
		defer gz.Close()

		c.Header("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		var e *gin.Engine
		e.ServeHTTP(gzipWriter{ResponseWriter: c.Writer, Writer: gz}, c.Request)
	}
}

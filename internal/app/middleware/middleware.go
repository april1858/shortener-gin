package middleware

import (
	"compress/gzip"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
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
		if c.GetHeader("Accept-Encoding") == "gzip" {
			log.Println("gzip")
			gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
			if err != nil {
				io.WriteString(c.Writer, err.Error())
				return
			}
			defer gz.Close()
			c.Writer = gzipWriter{ResponseWriter: c.Writer, Writer: gz}
			c.Next()
		}
		c.Next()
	}
}

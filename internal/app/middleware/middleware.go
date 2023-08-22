package middleware

import (
	"compress/gzip"
	"fmt"
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
	fmt.Println("MW1")
	return func(c *gin.Context) {
		fmt.Println("MV2")
		if c.GetHeader("Accept-Encoding") == "gzip" {
			fmt.Println("MV3")
			log.Println("gzip")
			gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
			if err != nil {
				io.WriteString(c.Writer, err.Error())
				return
			}
			defer gz.Close()
			fmt.Println("MV4")
			c.Writer = gzipWriter{ResponseWriter: c.Writer, Writer: gz}
			fmt.Println("MV5")
			c.Next()
		}
		fmt.Println("MV6")
		c.Next()
	}
}

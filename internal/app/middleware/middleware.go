package middleware

import (
	"io"

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
		c.Next()
		/*
			aE := c.GetHeader("Accept-Encoding")
			switch aE {
			case "":
				fmt.Println("none")
			case "gzip":
				gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
				if err != nil {
					io.WriteString(c.Writer, err.Error())
					return
				}
				defer gz.Close()
				c.Writer = gzipWriter{ResponseWriter: c.Writer, Writer: gz}
				c.Next()

			default:
				c.Next()
			}
		*/
	}
}

package middleware

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type MW struct {
}

func New() *MW {
	return &MW{}
}

func (mw *MW) GZIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			gzip.DefaultDecompressHandle(c)
		}
		fmt.Println("Accept[] - ", c.Accepted)
		if c.GetHeader("Accept-Encoding") == "gzip, deflate, br" {
			gzip.Gzip(gzip.DefaultCompression)
			fmt.Println("Accept[2] - ", c.Accepted)
		}
		fmt.Println("Accept[3] - ", c.Accepted)
		c.Next()
		fmt.Println("Accept[4] - ", c.Accepted)
	}
}

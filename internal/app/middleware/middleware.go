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

func (mv *MW) GZIP() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("Accept-Encoding") == "gzip" {
			fmt.Println("c1.Request - ", c.Request)
			fmt.Println("c1.Params - ", c.Params)
			if c.GetHeader("Content-Encoding") == "gzip" {
				fmt.Println("c2.Request - ", c.Request)
				fmt.Println("c2.Params - ", c.Params)
				gzip.DefaultDecompressHandle(c)
			}
			gzip.Gzip(gzip.DefaultCompression)
			fmt.Println("c3.Request - ", c.Request)
			fmt.Println("c3.Params - ", c.Params)
		}

		c.Next()
	}
}

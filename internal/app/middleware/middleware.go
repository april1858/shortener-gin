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
	fmt.Println("MW1")
	return func(c *gin.Context) {
		fmt.Println("MV2")
		if c.GetHeader("Accept-Encoding") == "gzip" {
			gzip.Gzip(gzip.DefaultCompression)
		}
		if c.GetHeader("Contenet-Encoding") == "gzip" {
			gzip.DefaultDecompressHandle(c)
		}
		fmt.Println("MV6")
		c.Next()
	}
}

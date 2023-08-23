package middleware

import (
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
			if c.GetHeader("Content-Encoding") == "gzip" {
				gzip.DefaultDecompressHandle(c)
			}
			gzip.Gzip(gzip.DefaultCompression)
		}
		c.Next()
	}
}

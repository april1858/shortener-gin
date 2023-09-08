package middleware

import (
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type MW struct {
	mx sync.RWMutex
}

func New() *MW {
	return &MW{}
}

func (mw *MW) GZIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			gzip.DefaultDecompressHandle(c)
		}
		if c.GetHeader("Accept-Encoding") == "gzip" {
			gzip.Gzip(gzip.DefaultCompression)
		}
		c.Next()
	}
}

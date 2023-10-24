package middleware

import (
	"strings"

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
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4", ".png", ".gif", ".jpeg", ".jpg", ".webp"}))
		}
		c.Next()
	}
}

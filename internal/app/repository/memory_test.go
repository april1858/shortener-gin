package repository

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkStore(b *testing.B) {
	r := NewMemStorage()
	var ctx *gin.Context
	b.Run("recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Store(ctx, "qqqqq", "aaaaa", "bbbbb")
		}
	})
}

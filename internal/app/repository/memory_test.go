package repository

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkStore(b *testing.B) {
	r := NewMemStorage()
	var ctx *gin.Context
	b.Run("store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Store(ctx, "http://abcd1234.ru", "1234abcd")
		}
	})

	b.Run("find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Find(ctx, "1234abcd")
		}
	})

	b.Run("findByUID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.FindByUID(ctx, "1234abcd")
		}
	})
}

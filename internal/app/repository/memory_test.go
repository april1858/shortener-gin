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
			r.Store(ctx, "a1234567", "http://a-a.ru", "123aaaa2")
		}
	})

	b.Run("find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Find(ctx, "a1234567")
		}
	})

	b.Run("findByUID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.FindByUID(ctx, "a1234567")
		}
	})
}

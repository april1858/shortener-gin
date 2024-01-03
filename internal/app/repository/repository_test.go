package repository

import (
	"github.com/gin-gonic/gin"
)

func ExampleMemory_StoreBatch() {
	batch := make([]map[string]string, 0)
	elem1 := make(map[string]string)
	elem1 = map[string]string{"short_url": "12345678", "original_url": "http://a-a.ru", "uid": "1"}
	batch[0] = elem1
	m := &Memory{}
	ctx := new(gin.Context)
	m.StoreBatch(ctx, batch)

	// output "12345678"
}

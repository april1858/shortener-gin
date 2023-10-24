package repository

import (
	"errors"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Memory struct {
	mx     sync.RWMutex
	memory []string
}

func NewMemStorage() *Memory {
	m := make([]string, 0, 1)
	return &Memory{memory: m}
}
func (r *Memory) Store(_ *gin.Context, short, original, uid string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.memory = append(r.memory, short+" "+original+" "+uid)
	return "", nil
}

func (r *Memory) Find(_ *gin.Context, short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range r.memory {
		var v = strings.Fields(value)
		if short == v[0] {
			return v[1], nil
		}
	}
	return "", nil
}

func (r *Memory) FindByUID(_ *gin.Context, uid string) ([]string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	answer := make([]string, 0, 4)
	for _, value := range r.memory {
		var v = strings.Fields(value)
		if uid == v[2] {
			answer = append(answer, v[0]+" "+v[1])
		}
	}
	if len(answer) == 0 {
		return nil, errors.New("NOT")
	}
	return answer, nil
}

func (r *Memory) Ping() (string, error) {
	return "Yes! Ping from Memory\n", nil
}

func (r *Memory) StoreBatch(_ *gin.Context, _ []map[string]string) error {
	return nil
}

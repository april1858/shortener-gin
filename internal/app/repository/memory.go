package repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type S struct {
	UID  string
	Data []string
}

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
	for _, vv := range f {
		if short == vv {
			return "deleted", nil
		}
	}
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

func (r *Memory) Delete(_ *gin.Context, c chan S) {
	go func() {
		//wg := &sync.WaitGroup{}
		var s = <-c
		data := s.Data
		fmt.Println("data - ", data)
		uid := s.UID
		fmt.Println("r.memory - ", r.memory)
		for _, rr := range data {
			f = append(f, rr)
			fmt.Println("f = ", f)
			for i, value := range r.memory {
				var v = strings.Fields(value)
				if uid == v[2] && rr == v[0] {
					copy(r.memory[i:], r.memory[i+1:])
					r.memory = r.memory[:len(r.memory)-1]
				}
			}

		}
		fmt.Println("r.memory2 - ", r.memory)
		//wg.Wait()
	}()
}

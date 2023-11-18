package repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	Store(ctx *gin.Context, short, originsl, uid string) (string, error)
	Find(ctx *gin.Context, short string) (string, error)
	FindByUID(*gin.Context, string) ([]string, error)
	StoreBatch(*gin.Context, []map[string]string) error
	Ping() (string, error)
	//Del(S)
}

type S struct {
	UID  string
	Data []string
}

var memory = make([]string, 0)

type Memory struct {
	mx     sync.RWMutex
	memory []string
}

var ch = make(chan S)

func New(c *config.Config) (Repository, chan S, error) {
	var r Repository
	var err error
	switch {
	case c.DatabaseDsn != "":
		r, err = NewDBStorage(c.DatabaseDsn)
		if err != nil {
			return nil, nil, err
		}
	case c.FileStoragePath != "":
		r = NewFileStorage(c.FileStoragePath)
	default:
		r = NewMemStorage()
	}

	return r, ch, nil
}

func NewMemStorage() *Memory {
	m := make([]string, 0, 1)
	p := &Memory{memory: m}
	go funnelm(p)
	fmt.Println("initm")
	return p
}
func (r *Memory) Store(_ *gin.Context, short, original, uid string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.memory = append(r.memory, short+" "+original+" "+uid+" "+"true")
	return "", nil
}

func (r *Memory) Find(_ *gin.Context, short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range r.memory {
		var v = strings.Fields(value)
		if short == v[0] {
			if v[3] == "false" {
				return "deleted", nil
			}
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

func funnelm(m *Memory) {
	fmt.Println("funnelm")
	v := <-ch
	data := v.Data
	uid := v.UID
	for _, rd := range data {
		for i, value := range m.memory {
			var v = strings.Fields(value)
			fmt.Println("rd - ", rd, "v[0] - ", v[0])
			if uid == v[2] && rd == v[0] {
				m.memory[i] = v[0] + " " + v[1] + " " + v[2] + " " + "false"
			}
		}
	}
	time.Sleep(time.Second * 60)
	Delm(m)
}

func Delm(m *Memory) {
	for i, value := range m.memory {
		var v = strings.Fields(value)
		if v[3] == "false" {
			m.memory = append(m.memory[:i], m.memory[i+1:]...)
		}
	}
}

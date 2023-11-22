package repository

import (
	"fmt"
	"sync"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/entity"
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

//var memory = make([]string, 0)

type eS entity.StoreElem

type Memory struct {
	mx     sync.RWMutex
	memory []eS
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
	m := make([]eS, 0, 1)
	p := &Memory{memory: m}
	go funnelm(p)
	return p
}
func (r *Memory) Store(_ *gin.Context, short, original, uid string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.memory = append(r.memory, eS{Short: short, Original: original, UID: uid, Condition: true})
	return "", nil
}

func (r *Memory) Find(_ *gin.Context, short string) (string, error) {
	fmt.Println("Find")
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range r.memory {
		if value.Short == short {
			if !value.Condition {
				return "", entity.ErrDeleted
			}
			return value.Original, nil
		}
	}
	return "", entity.ErrNotFound
}

func (r *Memory) FindByUID(_ *gin.Context, uid string) ([]string, error) {
	fmt.Println("FindByUID")
	r.mx.Lock()
	defer r.mx.Unlock()
	answer := make([]string, 0, 4)
	for _, value := range r.memory {
		if uid == value.UID {
			answer = append(answer, value.Short+" "+value.Original)
		}
	}
	if len(answer) == 0 {
		return nil, entity.ErrNotFound
	}
	return answer, nil
}

func (r *Memory) Ping() (string, error) {
	return "Yes! Ping from Memory\n", nil
}

func (r *Memory) StoreBatch(_ *gin.Context, batch []map[string]string) error {
	for _, v := range batch {
		r.memory = append(r.memory, eS{Short: v["short_url"], Original: v["original_url"], UID: v["uid"], Condition: true})
	}
	return nil
}

func funnelm(m *Memory) {
	for v := range ch {
		data := v.Data
		uid := v.UID
		for _, rd := range data {
			for i, value := range m.memory {
				if uid == value.UID && rd == value.Short {
					m.memory[i] = eS{Short: value.Short, Original: value.Original, UID: value.UID, Condition: false}
				}
			}
		}
	}
	Delm(m)
}

func Delm(m *Memory) {
	for i, value := range m.memory {
		if !value.Condition {
			m.memory = append(m.memory[:i], m.memory[i+1:]...)
		}
	}
}

package repository

//go:generate mockgen -build_flags=--mod=mod -destination mocks/postgres.go github.com/april1858/shortener-gin/internal/app/repository Repository
import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	_ "github.com/golang/mock/mockgen/model"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/entity"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	Store(ctx *gin.Context, originsl, uid string) (string, error)
	Find(ctx *gin.Context, short string) (string, error)
	FindByUID(*gin.Context, string) ([]string, error)
	StoreBatch(*gin.Context, []map[string]string) error
	Ping() (string, error)
	//Del(S)
}

type ES entity.StoreElem

type Memory struct {
	mx     sync.RWMutex
	Memory []ES
}

var ch = make(chan entity.ChData)

func New(c *config.Config) (Repository, chan entity.ChData, error) {
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
	m := make([]ES, 0, 1)
	p := &Memory{Memory: m}
	go funnelm(p)
	return p
}
func (r *Memory) Store(_ *gin.Context, original, uid string) (string, error) {
	short, err := GetRand()
	if err != nil {
		return "", err
	}
	r.mx.Lock()
	defer r.mx.Unlock()
	r.Memory = append(r.Memory, ES{Short: short, Original: original, UID: uid, Condition: true})
	return short, nil
}

func (r *Memory) Find(_ *gin.Context, short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range r.Memory {
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
	r.mx.Lock()
	defer r.mx.Unlock()
	answer := make([]string, 0, 4)
	for _, value := range r.Memory {
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
		r.Memory = append(r.Memory, ES{Short: v["short_url"], Original: v["original_url"], UID: v["uid"], Condition: true})
	}
	return nil
}

func funnelm(m *Memory) {
	for v := range ch {
		data := v.Data
		uid := v.UID
		for _, rd := range data {
			for i, value := range m.Memory {
				if uid == value.UID && rd == value.Short {
					m.Memory[i] = ES{Short: value.Short, Original: value.Original, UID: value.UID, Condition: false}
				}
			}
		}
	}
	Delm(m)
}

func Delm(m *Memory) {
	for i, value := range m.Memory {
		if !value.Condition {
			m.Memory = append(m.Memory[:i], m.Memory[i+1:]...)
		}
	}
}

// GetRand() creates a random number and return as a string - shortURL
func GetRand() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	short := hex.EncodeToString(b)
	return short, nil
}

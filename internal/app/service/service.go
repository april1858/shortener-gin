package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/april1858/shortener-gin/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	Store(ctx *gin.Context, short, originsl, uid string) (string, error)
	Find(ctx *gin.Context, short string) (string, error)
	FindByUID(*gin.Context, string) ([]string, error)
	StoreBatch(*gin.Context, []map[string]string) error
	Ping() (string, error)
	Del([]repository.S)
}

type Service struct {
	r Repository
}

var ch chan repository.S

func New(c *config.Config) (*Service, chan repository.S, error) {
	ch = make(chan repository.S)
	var r Repository
	var err error
	switch {
	case c.DatabaseDsn != "":
		r, err = repository.NewDBStorage(c.DatabaseDsn)
		if err != nil {
			return nil, nil, err
		}
	case c.FileStoragePath != "":
		r = repository.NewFileStorage(c.FileStoragePath)
	default:
		r = repository.NewMemStorage()
	}
	dests := Split(ch, 5)

	var wg sync.WaitGroup
	wg.Add(len(dests))
	for i, ch := range dests {
		go func(i int, d chan repository.S) {
			defer wg.Done()
			for val := range d {
				fmt.Printf("#%d got %v\n", i, val)
			}
		}(i, ch)
	}
	wg.Wait()

	return &Service{
		r: r,
	}, ch, nil
}

func (s *Service) CreatorShortened(ctx *gin.Context, originalURL string) (string, error) {
	uid := ctx.MustGet("UID").(string)
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error from CreatorShortened()", err
	}
	short := hex.EncodeToString(b)
	shorter, err := s.r.Store(ctx, short, originalURL, uid)
	if err != nil {
		return shorter, err
	}
	return short, nil
}

func (s *Service) FindOriginalURL(ctx *gin.Context, shortened string) (string, error) {
	answer, err := s.r.Find(ctx, shortened)
	return answer, err
}

func (s *Service) FindByUID(ctx *gin.Context) ([]string, error) {
	uid := ctx.MustGet("UID").(string)
	answer, err := s.r.FindByUID(ctx, uid)

	return answer, err
}

func (s *Service) CreatorShortenedBatch(ctx *gin.Context, batch []map[string]string) ([]string, error) {
	answer := make([]string, 0, 2)
	toDB := make([]map[string]string, 0)
	uid := ctx.MustGet("UID").(string)
	for _, v := range batch {
		mp := make(map[string]string, 0)
		b := make([]byte, 4)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}
		answer = append(answer, hex.EncodeToString(b)+" "+v["original_url"]+" "+uid)
		mp["short_url"] = hex.EncodeToString(b)
		mp["original_url"] = v["original_url"]
		mp["uid"] = uid
		toDB = append(toDB, mp)
	}
	err := s.r.StoreBatch(ctx, toDB)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *Service) Ping() (string, error) {
	answer, err := s.r.Ping()
	return answer, err
}

func Split(source chan repository.S, n int) []chan repository.S {
	dests := make([]chan repository.S, 0)
	for i := 0; i < n; i++ {
		// Создать n выходных каналов
		ch := make(chan repository.S)
		dests = append(dests, ch)
		go func() {
			defer close(ch)
			for val := range source {
				ch <- val
			}
		}()
	}
	return dests
}

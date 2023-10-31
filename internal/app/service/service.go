package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

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
	Delete(*gin.Context, chan repository.S)
}

type Service struct {
	r Repository
}

func New(c *config.Config) (*Service, error) {
	var r Repository
	var err error
	switch {
	case c.DatabaseDsn != "":
		r, err = repository.NewDBStorage(c.DatabaseDsn)
		if err != nil {
			return nil, err
		}
	case c.FileStoragePath != "":
		r = repository.NewFileStorage(c.FileStoragePath)
	default:
		r = repository.NewMemStorage()
	}
	return &Service{
		r: r,
	}, nil
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

func (s *Service) Delete(ctx *gin.Context, c chan repository.S) {
	fmt.Println("ccc - ", c)
	s.r.Delete(ctx, c)
}

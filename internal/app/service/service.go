package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	Store(ctx *gin.Context, short, originsl string) error
	Find(ctx *gin.Context, short string) (string, error)
	FindByUID(ctx *gin.Context) ([]string, error)
	StoreBatch(string, []map[string]string) error
	Ping(dsn string) (string, error)
}

type Service struct {
	r Repository
	c config.Config
}

func New(r Repository, c config.Config) *Service {
	return &Service{
		r: r,
		c: c,
	}
}

func (s *Service) CreatorShortened(ctx *gin.Context, originalURL string) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in CreatorShortened()", err
	}
	short := hex.EncodeToString(b)
	err = s.r.Store(ctx, short, originalURL)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *Service) CreatorShortenedBatch(batch []map[string]string, uid string) []string {
	answer := make([]string, 0, 2)
	toDB := make([]map[string]string, 0)

	for _, v := range batch {
		mp := make(map[string]string, 0)
		b := make([]byte, 4)
		_, err := rand.Read(b)
		if err != nil {
			return nil
		}
		answer = append(answer, hex.EncodeToString(b)+" "+v["original_url"]+" "+uid)
		mp["short_url"] = hex.EncodeToString(b)
		mp["original_url"] = v["original_url"]
		mp["uid"] = uid
		toDB = append(toDB, mp)
	}
	err := s.r.BulkInsert(s.c.DatabaseDsn, toDB)
	if err != nil {
		fmt.Println("err from service - ", err)
	}
	return answer
}

func (s *Service) FindOriginalURL(ctx *gin.Context, shortened string) (string, error) {
	answer, err := s.r.Find(ctx, shortened)
	return answer, err
}

func (s *Service) FindByUID(ctx *gin.Context) ([]string, error) {
	answer, err := s.r.FindByUID(ctx)

	return answer, err
}

func (s *Service) Ping() (string, error) {
	answer, err := s.r.Ping(s.c.DatabaseDsn)

	return answer, err
}

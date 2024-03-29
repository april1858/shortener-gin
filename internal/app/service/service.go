package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/april1858/shortener-gin/internal/app/entity"
	"github.com/april1858/shortener-gin/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Service struct {
	r repository.Repository
}

func New(r repository.Repository, ch chan entity.ChData) (*Service, chan entity.ChData) {
	return &Service{
		r: r,
	}, ch
}

func (s *Service) CreatorShortened(ctx *gin.Context, originalURL string) (string, error) {
	uid := ctx.MustGet("UID").(string)
	short, err := s.r.Store(ctx, originalURL, uid)
	if err != nil {
		return short, err
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
		//answer = append(answer, hex.EncodeToString(b)+" "+v["original_url"]+" "+uid)
		answer = append(answer, hex.EncodeToString(b))
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

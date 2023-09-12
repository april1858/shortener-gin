package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
)

type Repository interface {
	Store(short, original, uid, string) error
	Find(short string) (string, error)
	FindByUID(uid string) ([]string, error)
	StoreBatch(string, []map[string]string) error
	Ping(dsn string) (string, error)
}

type Service struct {
	r Repository
}

func New(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) CreatorShortened(originalURL, uid string) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in CreatorShortened()", err
	}
	s.r.Store(hex.EncodeToString(b), originalURL, uid)
	return hex.EncodeToString(b), nil
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
	err := s.r.StoreBatch(s.c.DatabaseDsn, toDB)
	if err != nil {
		fmt.Println("err from service - ", err)
	}
	return answer
}

func (s *Service) FindOriginalURL(shortened string) (string, error) {
	var (
		answer string
		err    error
	)
	answer, err = s.r.Find(shortened)
	return answer, err
}

func (s *Service) FindByUID(uid string) ([]string, error) {
	var (
		answer []string
		err    error
	)
	answer, err = s.r.FindByUID(uid)
	return answer, err
}

func (s *Service) Ping() (string, error) {
	answer, err := s.r.Ping()

	return answer, err
}

package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
)

type Repository interface {
	MemoryStore(short, original, uid string) error
	MemoryFind(short string) (string, error)
	MemoryFindByUID(uid string) ([]string, error)
	FileStore(filename, short, original, uid string) error
	FileFind(filename, short string) (string, error)
	FileFindByUID(filename, uid string) ([]string, error)
	DBStore(dsn, short, original, uid string) (string, error)
	DBFind(dsn, shorturl string) (string, error)
	DBFindByUID(dsn, uid string) ([]string, error)
	Ping(dsn string) (string, error)
	BulkInsert(string, []map[string]string) error
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

func (s *Service) CreatorShortened(originalURL, uid string) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in CreatorShortened()", err
	}
	switch {
	case s.c.FileStoragePath != "":
		s.r.FileStore(s.c.FileStoragePath, hex.EncodeToString(b), originalURL, uid)
	case s.c.DatabaseDsn != "":
		answer, err := s.r.DBStore(s.c.DatabaseDsn, hex.EncodeToString(b), originalURL, uid)
		if err != nil {
			return answer, err
		}
	default:
		s.r.MemoryStore(hex.EncodeToString(b), originalURL, uid)
	}

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
	err := s.r.BulkInsert(s.c.DatabaseDsn, toDB)
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
	switch {
	case s.c.FileStoragePath != "":
		answer, err = s.r.FileFind(s.c.FileStoragePath, shortened)
	case s.c.DatabaseDsn != "":
		answer, err = s.r.DBFind(s.c.DatabaseDsn, shortened)
	default:
		answer, err = s.r.MemoryFind(shortened)
	}
	return answer, err
}

func (s *Service) FindByUID(uid string) ([]string, error) {
	var (
		answer []string
		err    error
	)
	switch {
	case s.c.FileStoragePath != "":
		answer, err = s.r.FileFindByUID(s.c.FileStoragePath, uid)
	case s.c.DatabaseDsn != "":
		answer, err = s.r.DBFindByUID(s.c.DatabaseDsn, uid)
	default:
		answer, err = s.r.MemoryFindByUID(uid)
	}

	return answer, err
}

func (s *Service) Ping() (string, error) {
	answer, err := s.r.Ping(s.c.DatabaseDsn)

	return answer, err
}

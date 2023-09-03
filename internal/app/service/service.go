package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/april1858/shortener-gin/internal/app/config"
)

var UID string

type Repository interface {
	MemoryStore(short, original string) error
	MemoryFind(short string) (string, error)
	MemoryFindByUID() ([]string, error)
	FileStore(filename, short, original string) error
	FileFind(filename, short string) (string, error)
	FileFindByUID(filename string) ([]string, error)
	DBStore(dsn, short, original string) error
	DBFind(dsn, shorturl string) (string, error)
	DBFindByUID(dsn string) ([]string, error)
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

func (s *Service) CreatorShortened(originalURL string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in CreatorShortened()"
	}
	switch {
	case s.c.FileStoragePath != "":
		s.r.FileStore(s.c.FileStoragePath, hex.EncodeToString(b), originalURL)
	case s.c.DatabaseDsn != "":
		s.r.DBStore(s.c.DatabaseDsn, hex.EncodeToString(b), originalURL)
	default:
		s.r.MemoryStore(hex.EncodeToString(b), originalURL)
	}

	return hex.EncodeToString(b)
}

func (s *Service) CreatorShortenedBatch(batch []map[string]string) []string {
	answer := make([]string, 0, 2)
	toDB := make([]map[string]string, 2, 2)

	for i, v := range batch {
		mp := make(map[string]string, 0)
		b := make([]byte, 4)
		_, err := rand.Read(b)
		if err != nil {
			return nil
		}
		answer = append(answer, hex.EncodeToString(b)+" "+v["original_url"]+" "+UID)
		mp["short_url"] = hex.EncodeToString(b)
		mp["original_url"] = v["original_url"]
		mp["uid"] = UID
		toDB[i] = mp
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

func (s *Service) FindByUID() ([]string, error) {
	var (
		answer []string
		err    error
	)
	switch {
	case s.c.FileStoragePath != "":
		answer, err = s.r.FileFindByUID(s.c.FileStoragePath)
	case s.c.DatabaseDsn != "":
		answer, err = s.r.DBFindByUID(s.c.DatabaseDsn)
	default:
		answer, err = s.r.MemoryFindByUID()
	}

	return answer, err
}

func (s *Service) Ping() (string, error) {
	answer, err := s.r.Ping(s.c.DatabaseDsn)

	return answer, err
}

package service

import (
	"crypto/rand"
	"encoding/hex"
)

type Repository interface {
	Store(string, string) error
	Find(string) (string, error)
	FindAllUID() ([]string, error)
}

type Service struct {
	R Repository
}

func New(r Repository) *Service {
	return &Service{
		R: r,
	}
}

func (s *Service) CreatorShortened(originalURL string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in CreatorShortened()"
	}

	s.R.Store(hex.EncodeToString(b), originalURL)

	return hex.EncodeToString(b)
}

func (s *Service) FindOriginalURL(shortened string) (string, error) {
	answer, err := s.R.Find(shortened)

	return answer, err
}

func (s *Service) FindAllUID() ([]string, error) {
	answer, err := s.R.FindAllUID()

	return answer, err
}

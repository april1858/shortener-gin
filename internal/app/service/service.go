package service

import (
	"crypto/rand"
	"encoding/hex"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) CreatorShortened() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "error in createCode()"
	}

	return hex.EncodeToString(b)
}

package repository

import (
	"errors"
	"strings"
	"sync"
)

var M = make([]string, 0, 1)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) MemoryStore(short, original, uid string) error {
	mx := new(sync.RWMutex)
	mx.Lock()
	defer mx.Unlock()
	M = append(M, short+" "+original+" "+uid)
	return nil
}

func (r *Repository) MemoryFind(short string) (string, error) {
	for _, value := range M {
		var v = strings.Fields(value)
		if short == v[0] {
			return v[1], nil
		}
	}
	return "", nil
}

func (r *Repository) MemoryFindByUID(uid string) ([]string, error) {
	answer := make([]string, 0, 4)
	for _, value := range M {
		var v = strings.Fields(value)
		if uid == v[2] {
			answer = append(answer, v[0]+" "+v[1])
		}
	}
	if len(answer) == 0 {
		return nil, errors.New("NOT")
	}
	return answer, nil
}

package repository

import (
	"errors"
	"sync"
)

type Repository struct {
	mx     *sync.RWMutex
	M map[string]string
}

func New() *Repository {
	m := make(map[string]string)
	mx := new(sync.RWMutex)
	return &Repository{M: m, mx: mx}
}

func (r *Repository) Store(short, original string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.M[short] = original
}

func (r *Repository) Find(short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.M[short]
	if !ok {
		return "", errors.New("Нет соответствия ...")
	}else{
		return r.M[short], nil
	}
}

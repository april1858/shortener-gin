package repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var M = make([]string, 0, 1)
var UID string

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) MemoryStore(short, original string) error {
	mx := new(sync.RWMutex)
	mx.Lock()
	defer mx.Unlock()
	M = append(M, short+" "+original+" "+UID)
	fmt.Println("M - ", M)
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

func (r *Repository) MemoryFindByUID() ([]string, error) {
	answer := make([]string, 0, 4)
	for _, value := range M {
		var v = strings.Fields(value)
		if UID == v[2] {
			answer = append(answer, v[0]+" "+v[1])
		}
	}
	if len(answer) == 0 {
		return nil, errors.New("NOT")
	}
	return answer, nil
}

package repository

import (
	"errors"
	"strings"
	"sync"
)

type Memory struct {
	mx     sync.RWMutex
	memory []string
}

func NewMemStorage() *Memory {
	m := make([]string, 0, 1)
	return &Memory{memory: m}
}

/*
	func (r *Repository) Store(ctx *gin.Context, short, original string) (string, error) {
		var err error
		uid := ctx.MustGet("UID").(string)
		switch {
		case config.Cnf.FileStoragePath != "":
			err = r.FileStore(config.Cnf.FileStoragePath, short, original, uid)
			if err != nil {
				return "", err
			}
		case config.Cnf.DatabaseDsn != "":
			isShort, err := r.PGSStore(ctx, short, original, uid)
			if err != nil {
				return isShort, err
			}
		default:
			err = r.MemoryStore(short, original, uid)
			if err != nil {
				return "", err
			}
		}
		return "", nil
	}

	func (r *Repository) Find(ctx *gin.Context, short string) (string, error) {
		var answer string
		var err error
		switch {
		case config.Cnf.FileStoragePath != "":
			answer, err = r.FileFind(config.Cnf.FileStoragePath, short)
		case config.Cnf.DatabaseDsn != "":
			answer, err = r.PGSFind(ctx, short)
		default:
			answer, err = r.MemoryFind(short)
		}
		return answer, err
	}

	func (r *Repository) FindByUID(ctx *gin.Context) ([]string, error) {
		var answer []string
		var err error
		uid := ctx.MustGet("UID").(string)
		switch {
		case config.Cnf.FileStoragePath != "":
			answer, err = r.FileFindByUID(config.Cnf.FileStoragePath, uid)
		case config.Cnf.DatabaseDsn != "":
			answer, err = r.PGSFindByUID(ctx, uid)
		default:
			answer, err = r.MemoryFindByUID(uid)
		}
		return answer, err
	}

	func (r *Repository) StoreBatch(ctx *gin.Context, batch []map[string]string) error {
		var err error
		switch {
		case config.Cnf.FileStoragePath != "":
			err = errors.New("pass")
		case config.Cnf.DatabaseDsn != "":
			err = r.PGSStoreBatch(ctx, batch)
		default:
			err = errors.New("pass")
		}
		return err
	}
*/
func (r *Memory) Store(short, original, uid string) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.memory = append(r.memory, short+" "+original+" "+uid)
	return nil
}

func (r *Memory) Find(short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range r.memory {
		var v = strings.Fields(value)
		if short == v[0] {
			return v[1], nil
		}
	}
	return "", nil
}

func (r *Memory) FindByUID(uid string) ([]string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	answer := make([]string, 0, 4)
	for _, value := range r.memory {
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

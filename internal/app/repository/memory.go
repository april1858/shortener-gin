package repository

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var M = make([]string, 0, 1)

type Repository struct {
	ctx context.Context
	db  *pgxpool.Pool
	mx  sync.RWMutex
}

func New(cnf *config.Config) *Repository {
	if cnf.DatabaseDsn != "" {
		ctx := context.Background()
		poolConfig, err := pgxpool.ParseConfig(cnf.DatabaseDsn)
		if err != nil {
			return nil
		}
		db, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			return nil
		}
		_, err = db.Exec(ctx, `create table if not exists shortener ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE)`)
		if err != nil {
			return nil
		}

		return &Repository{
			ctx: ctx,
			db:  db,
		}
	}
	return &Repository{
		ctx: nil,
		db:  nil,
	}
}

func (r *Repository) Store(short, original, uid string) error {
	switch {
	case s.c.FileStoragePath != "":
		answer, err = r.FileStore(s.c.FileStoragePath, shortened)
	case s.c.DatabaseDsn != "":
		answer, err = r.DBStore(s.c.DatabaseDsn, shortened)
	default:
		answer, err = r.MemoryStore(short, original, uid)
	}
	return answer, err
}

func (r *Repository) Find(short string) (string, error) {
	switch {
	case s.c.FileStoragePath != "":
		answer, err = r.FileFind(short)
	case s.c.DatabaseDsn != "":
		answer, err = r.DBFind(short)
	default:
		answer, err = r.MemoryFind(short
		)
	}
	return answer, err
}

func (r *Repository) MemoryStore(short, original, uid string) error {

	r.mx.Lock()
	defer r.mx.Unlock()
	M = append(M, short+" "+original+" "+uid)
	return nil
}

func (r *Repository) MemoryFind(short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range M {
		var v = strings.Fields(value)
		if short == v[0] {
			return v[1], nil
		}
	}
	return "", nil
}

func (r *Repository) MemoryFindByUID(uid string) ([]string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
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

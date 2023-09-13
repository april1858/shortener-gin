package repository

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/april1858/shortener-gin/internal/app/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var memory = make([]string, 0, 1)

type Repository struct {
	mx      sync.RWMutex
	connPGS *pgxpool.Pool
}

func New() *Repository {
	if config.Cnf.DatabaseDsn != "" {
		var ctx *gin.Context
		var db *pgxpool.Pool
		poolConfig, err := pgxpool.ParseConfig(config.Cnf.DatabaseDsn)
		if err != nil {
			log.Fatalln("Unable to parse DATABASE_DSN:", err)
		}
		db, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			log.Fatalln("Unable to create connection pool:", err)
		}
		_, err = db.Exec(ctx, `create table if not exists shortener ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE)`)
		if err != nil {
			log.Fatal("Not create table - ", err)
		}
		return &Repository{
			connPGS: db,
		}
	}

	return &Repository{}
}

func (r *Repository) Ping(ctx *gin.Context) (string, error) {
	message, err := r.PGSPing(ctx, config.Cnf.DatabaseDsn)
	if err != nil {
		return "", err
	}
	return message, nil
}

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
		isShort, err := r.PGSStore(ctx, config.Cnf.DatabaseDsn, short, original, uid)
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

func (r *Repository) Find(short string) (string, error) {
	switch {
	case config.Cnf.FileStoragePath != "":
		answer, err = r.FileFind(short)
	case config.Cnf.DatabaseDsn != "":
		answer, err = r.PGSFind(short)
	default:
		answer, err = r.MemoryFind(short)
	}
	return answer, err
}

func (r *Repository) MemoryStore(short, original, uid string) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	memory = append(memory, short+" "+original+" "+uid)
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

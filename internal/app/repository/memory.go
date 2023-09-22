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
	cnf string
}

func New(cnf config.Config) *Repository {
	if cnf.DatabaseDsn != "" {
		ctx := new(gin.Context)
		var db *pgxpool.Pool
		poolConfig, err := pgxpool.ParseConfig(cnf.DatabaseDsn)
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
			cnf: cnf.DatabaseDsn,
		}
	}
	return &Repository{}
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

func (r *Repository) MemoryStore(short, original, uid string) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	memory = append(memory, short+" "+original+" "+uid)
	return nil
}

func (r *Repository) MemoryFind(short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, value := range memory {
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
	for _, value := range memory {
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

package repository

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//_, err = db.Exec(ctx, `create table if not exists shortener ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE)`)

func (r *Repository) PGSPing(ctx *gin.Context, dsn string) (string, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)
	err = conn.Ping(ctx)
	if err != nil {
		fmt.Println("Not connecting!")
		return "", err
	} else {
		fmt.Println("Yes, connecting!")
	}
	return "Conn", nil
}

func (r *Repository) PGSStore(ctx *gin.Context, short, original, uid string) (string, error) {
	db := r.connPGS
	if _, err := db.Exec(ctx, `insert into "shortener" (uid, short_url, original_url) values ($1,$2,$3)`, uid, short, original); err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				var answer string
				row := db.QueryRow(ctx, `select short_url from "shortener" where original_url=$1`, original)
				err := row.Scan(&answer)
				if err != nil {
					panic(err)
				}
				return answer, pgxError
			}
		}
	}
	return "", nil
}

func (r *Repository) StoreBatch(ctx *gin.Context, batch []map[string]string) error {
	db := r.connPGS
	_, err := db.Exec(ctx, `INSERT INTO shortener (uid, short_url, original_url) VALUES ($1, $2, $3)`, batch)
	if err != nil {
		fmt.Println("222 - ", err)
	}
	return nil
}

func (r *Repository) PGSFind(ctx *gin.Context, short string) (string, error) {
	var answer string
	db := r.connPGS
	row := db.QueryRow(ctx, `select original_url from "shortener" where short_url=$1`, short)
	err := row.Scan(&answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (r *Repository) DBFindByUID(ctx *gin.Context, uid string) ([]string, error) {
	var answer []string
	db := r.connPGS
	row := db.QueryRow(ctx, `select original_url from "shortener" where uid=uid`)
	err := row.Scan(&answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *Repository) BulkInsert(ctx *gin.Context, bulks []map[string]string) error {
	db := r.connPGS
	query := `INSERT INTO shortener (uid, short_url, original_url) VALUES (@uid, @short_url, @original_url)`
	batch := &pgx.Batch{}
	for _, bulk := range bulks {
		args := pgx.NamedArgs{
			"uid":          bulk["uid"],
			"short_url":    bulk["short_url"],
			"original_url": bulk["original_url"],
		}
		batch.Queue(query, args)
	}
	results := db.SendBatch(ctx, batch)
	defer results.Close()
	return nil
}

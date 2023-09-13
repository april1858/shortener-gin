package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)

//_, err = db.Exec(ctx, `create table if not exists shortener ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE)`)

func (r *Repository) PGSPing(ctx *gin.Context, dsn string) (string, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close(cx)
	err = conn.Ping(cx)
	if err != nil {
		fmt.Println("Not connecting!")
		return "", err
	} else {
		fmt.Println("Yes, connecting!")
	}
	return "Conn", nil
}

func (r *Repository) connectDB(ctx *gin.Context, dsn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_DSN:", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}
	_, err = db.Exec(ctx, `create table if not exists shortener ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE)`)
	if err != nil {
		return nil, nil, err
	}
	return db, nil
}

func (r *Repository) PGSStore(ctx *gin.Context,dsn, short, original, uid string) (string, error) {
	db, err := r.connectDB(ctx, dsn)
	if err != nil {
		return "", err
	}
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

func (r *Repository) StoreBatch(dsn string, batch []map[string]string) error {
	ctx, db, err := r.connectDB(dsn)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, `INSERT INTO shortener (uid, short_url, original_url) VALUES ($1, $2, $3)`, batch)
	if err != nil {
		fmt.Println("222 - ", err)
	}
	return nil
}

func (r *Repository) PGSFind(dsn, shorturl string) (string, error) {
	var answer string
	ctx, db, err := r.connectDB(dsn)
	if err != nil {
		return "", err
	}
	row := db.QueryRow(ctx, `select original_url from "shortener" where short_url=$1`, shorturl)
	err = row.Scan(&answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (r *Repository) DBFindByUID(dsn, uid string) ([]string, error) {
	var answer []string
	ctx, db, err := r.connectDB(dsn)
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(ctx, `select original_url from "shortener" where uid=uid`)
	err = row.Scan(&answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *Repository) BulkInsert(dsn string, bulks []map[string]string) error {
	ctx, db, err := r.connectDB(dsn)
	if err != nil {
		return err
	}
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

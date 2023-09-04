package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (r Repository) Ping(dsn string) (string, error) {
	cx := context.Background()
	conn, err := pgx.Connect(cx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close(cx)
	err = conn.Ping(cx)
	if err != nil {
		fmt.Println("Panic")
		return "", err
	} else {
		fmt.Println("Yes, connecting!")
	}
	return "Conn", nil
}

func (r Repository) connectDB(dsn string) (context.Context, *pgxpool.Pool) {
	ctx := context.Background()
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
		fmt.Println("err - ", err)
	}
	return ctx, db
}

func (r Repository) DBStore(dsn, short, original string) (string, error) {
	uid := UID
	ctx, db := r.connectDB(dsn)
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
				fmt.Println("answer 23505", answer)
				return answer, pgxError
			}
		}
	}
	return "", nil
}

func (r Repository) StoreBatch(dsn string, batch []map[string]string) error {
	ctx, db := r.connectDB(dsn)
	_, err := db.Exec(ctx, `INSERT INTO shortener (uid, short_url, original_url) VALUES ($1, $2, $3)`, batch)
	if err != nil {
		fmt.Println("222 - ", err)
	}
	return nil
}

func (r Repository) DBFind(dsn, shorturl string) (string, error) {
	var answer string
	ctx, db := r.connectDB(dsn)
	row := db.QueryRow(ctx, `select original_url from "shortener" where short_url=$1`, shorturl)
	err := row.Scan(&answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (r Repository) DBFindByUID(dsn string) ([]string, error) {
	var answer []string
	ctx, db := r.connectDB(dsn)
	row := db.QueryRow(ctx, `select original_url from "shortener" where uid=UID`)
	err := row.Scan(&answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (r Repository) BulkInsert(dsn string, bulks []map[string]string) error {
	ctx, db := r.connectDB(dsn)
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

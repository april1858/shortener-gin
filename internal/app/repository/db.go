package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (r Repository) Ping() (string, error) {
	cx := context.Background()
	conn, err := pgx.Connect(cx, r.c.DatabaseDsn)
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

func (r Repository) connectDB() (context.Context, *pgxpool.Pool) {
	ctx := context.Background()
	poolConfig, err := pgxpool.ParseConfig(os.Getenv(r.c.DatabaseDsn))
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

func (r Repository) insertDB(short, original string) error {
	ctx, db := r.connectDB()
	if _, err := db.Exec(ctx, `insert into "shortener" (uid, short_url, original_url) values ($1,$2,$3)`, short, original); err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				var answer string
				row := db.QueryRow(ctx, `select code from "shortener4" where url=$1`, original)
				err := row.Scan(&answer)
				if err != nil {
					panic(err)
				}
				return nil
			}
		}
	}
	return nil
}

/*
func (r Repository) StoreLot(redirect *entity.Redirect2) error {
	ctx, db := r.connectDB()
	if _, err := db.Exec(ctx, `insert into "shortener4" (str_id, code, url) values ($1,$2,$3)`, redirect.CorrelationID, redirect.ShortURL, redirect.OriginalURL); err != nil {
		fmt.Println("errorrs!", err)
	}
	r.memory[redirect.ShortURL] = redirect.OriginalURL
	return nil
}

func (mr Repository) findDB(shorturl string) string {
	var answer string
	ctx, db := connectDB()
	row := db.QueryRow(ctx, `select url from "shortener4" where code=$1`, shorturl)
	err := row.Scan(&answer)
	if err != nil {
		panic(err)
	}
	return answer
}
*/

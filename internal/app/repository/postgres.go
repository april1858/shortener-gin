package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	connPGS *pgxpool.Pool
	db      string
}

func NewDBStorage(db string) (*DB, error) {
	ctx := new(gin.Context)
	var conn *pgxpool.Pool
	poolConfig, err := pgxpool.ParseConfig(db)
	if err != nil {
		return nil, err
	}
	conn, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(ctx, `create table if not exists shortener6 ("id" SERIAL PRIMARY KEY, "uid" varchar(100), "short_url" varchar(50), "original_url" text UNIQUE, "condition" boolean DEFAULT true)`)
	if err != nil {
		return nil, err
	}
	return &DB{
		connPGS: conn,
		db:      db,
	}, nil
}

func (d *DB) Ping() (string, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, d.db)
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

func (d *DB) Store(ctx *gin.Context, short, original, uid string) (string, error) {
	db := d.connPGS
	if _, err := db.Exec(ctx, `insert into "shortener6" (uid, short_url, original_url) values ($1,$2,$3)`, uid, short, original); err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				var answer string
				row := db.QueryRow(ctx, `select short_url from "shortener6" where original_url=$1`, original)
				err := row.Scan(&answer)
				if err != nil {
					return "", err
				}
				return answer, pgxError
			}
		}
	}
	return "", nil
}

func (d *DB) Find(ctx *gin.Context, short string) (string, error) {
	var answer string
	db := d.connPGS
	row := db.QueryRow(ctx, `select original_url from shortener6 where short_url=$1`, short)
	err := row.Scan(&answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (d *DB) FindByUID(ctx *gin.Context, uid string) ([]string, error) {
	answer := make([]string, 0, 1)
	db := d.connPGS
	rows, err := db.Query(ctx, `select short_url, original_url from shortener6 where uid=$1`, uid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a1 string
		var a2 string
		err = rows.Scan(&a1, &a2)
		answer = append(answer, a1+" "+a2)
	}
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (d *DB) StoreBatch(ctx *gin.Context, bulks []map[string]string) error {
	db := d.connPGS
	query := `INSERT INTO shortener6 (uid, short_url, original_url) VALUES (@uid, @short_url, @original_url)`
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
	return results.Close()
}

func (d *DB) Delete(ctx *gin.Context, remove []string) (int64, error) {
	uid := ctx.MustGet("UID").(string)
	db := d.connPGS
	var removed int64
	for _, r := range remove {
		x, err := db.Exec(ctx, `UPDATE "shortener6" SET condition = false WHERE uid = $1 AND short_url = $2`, uid, r)
		removed = x.RowsAffected()
		if err != nil {
			return 0, err
		}
	}
	return removed, nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

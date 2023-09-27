package config

import (
	"flag"
	"os"
)

type Config struct {
	BaseURL         string
	ServerAddress   string
	FileStoragePath string
	DatabaseDsn     string
}

var (
	A = flag.String("a", "", "server_address")
	B = flag.String("b", "", "base_url")
	F = flag.String("f", "", "file_storage_path")
	D = flag.String("d", "", "database_dsn")
)

var Cnf Config

func New() *Config {

	var address, baseurl, file, db string
	a := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	f := os.Getenv("FILE_STORAGE_PATH")
	d := os.Getenv("DATABASE_DSN")
	flag.Parse()

	if a == "" {
		if *A == "" {
			address = "localhost:8080"
		} else {
			address = *A
		}
	} else {
		address = a
	}

	if b == "" {
		if *B == "" {
			baseurl = "http://" + address
		} else {
			baseurl = *B
		}
	} else {
		baseurl = b
	}

	if d == "" {
		db = *D
	} else {
		db = d
	}

	if f == "" {
		file = *F
	} else {
		file = f
	}

	Cnf.BaseURL = baseurl + "/"
	Cnf.ServerAddress = address
	Cnf.FileStoragePath = file
	Cnf.DatabaseDsn = db

	return &Cnf
}

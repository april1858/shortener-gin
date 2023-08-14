package config

import (
	"flag"
	"os"
)

type Config struct {
	BaseURL         string
	ServerAddress   string
	FileStoragePath string
	//DatabaseDsn	string
}

var (
	A = flag.String("a", "", "server_address")
	B = flag.String("b", "", "base_url")
	F = flag.String("f", "", "file_storage_path")
	//D = flag.String("d", "", "database_dsn")
)

var cnf Config

func New() *Config {

	var address, baseurl, file string
	a := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	f := os.Getenv("FILE_STORAGE_PATH")
	//d := os.Getenv("DATABASE_DSN")
	flag.Parse()

	if a == "" {
		if *A == "" {
			address = "8080"
		} else {
			address = *A
		}
	} else {
		address = a
	}

	if b == "" {
		if *B == "" {
			baseurl = "http://localhost" + ":" + address
		} else {
			baseurl = *B
		}
	} else {
		baseurl = b
	}

	if f == "" {
		if *F == "" {
			file = ""
		} else {
			file = *F
		}
	} else {
		file = f
	}

	cnf.BaseURL = baseurl
	cnf.ServerAddress = address
	cnf.FileStoragePath = file

	return &cnf
}
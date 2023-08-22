package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/april1858/shortener-gin/internal/app/config"
)

var M = make([]string, 0, 1)

type Repository struct {
	mx *sync.RWMutex
	c  *config.Config
}

func New(c *config.Config) *Repository {
	mx := new(sync.RWMutex)
	return &Repository{mx: mx, c: c}
}

func (r *Repository) Store(short, original string) error {
	fmt.Println("from Store")
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.c.FileStoragePath == "" {
		M = append(M, short+" "+original)
	} else {
		filename := r.c.FileStoragePath
		_, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				os.OpenFile(filename, os.O_CREATE, 0777)
			}
			M = append(M, short+" "+original)
		} else {
			content, err := os.ReadFile(filename)
			if err != nil {
				log.Println("error - ", err)
				return err
			}
			json.Unmarshal(content, &M)
			M = append(M, short+" "+original)
		}

		data, err := json.Marshal(M)
		if err != nil {
			log.Println("error ", err)
			return err
		}
		err = os.WriteFile(filename, data, 0777)
		if err != nil {
			log.Println("error ", err)
			return err
		}
	}
	return nil
}

func (r *Repository) Find(short string) (string, error) {
	if r.c.FileStoragePath == "" {
		for _, value := range M {
			var v = strings.Fields(value)
			if short == v[0] {
				return v[1], nil
			}
		}
	} else {
		filename := r.c.FileStoragePath
		fileData, err := os.ReadFile(filename)

		if err != nil {
			log.Println("error ", err)
			return "", err
		}
		parseData := []string{}
		json.Unmarshal(fileData, &parseData)

		for _, value := range parseData {
			var v = strings.Fields(value)
			if short == v[0] {
				return v[1], nil
			}
		}
	}
	return "", nil
}

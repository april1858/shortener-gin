package repository

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

type File struct {
	mx   *sync.RWMutex
	file string
}

func NewFile(f string) *File {
	mx := new(sync.RWMutex)
	return &File{file: f, mx: mx}
}

func (f File) Store(short, original string) error {

	f.mx.Lock()
	defer f.mx.Unlock()
	sm := make([]string, 0, 2)
	filename := f.file
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			os.OpenFile(filename, os.O_CREATE, 0777)
		}
		sm = append(sm, short+" "+original)
	} else {
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Println("error - ", err)
			return err
		}
		json.Unmarshal(content, &sm)
		sm = append(sm, short+" "+original)
	}

	data, err := json.Marshal(sm)
	if err != nil {
		log.Println("error ", err)
		return err
	}
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		log.Println("error ", err)
		return err
	}
	return nil
}

func (f File) Find(uid string) (string, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	filename := f.file
	fileData, err := os.ReadFile(filename)

	if err != nil {
		log.Println("error ", err)
		return "", err
	}
	parseData := []string{}
	json.Unmarshal(fileData, &parseData)

	for _, value := range parseData {
		var v = strings.Fields(value)
		if uid == v[0] {
			return v[1], nil
		}
	}
	return "", nil
}

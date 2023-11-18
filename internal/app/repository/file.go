package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type File struct {
	mx       sync.RWMutex
	filename string
}

func NewFileStorage(f string) *File {
	fmt.Println("initf")
	go funnelf()
	return &File{filename: f}
}

func (f *File) Store(_ *gin.Context, short, original, uid string) (string, error) {
	data := make([]string, 0, 1)
	f.mx.Lock()
	defer f.mx.Unlock()
	_, err := os.Stat(f.filename)
	if err != nil {
		if os.IsNotExist(err) {
			os.OpenFile(f.filename, os.O_CREATE, 0777)
		}
		data = append(data, short+" "+original+" "+uid)
	} else {
		content, err := os.ReadFile(f.filename)
		if err != nil {
			log.Println("error - ", err)
			return "", err
		}
		json.Unmarshal(content, &data)
		data = append(data, short+" "+original+" "+uid)
	}

	out, err := json.Marshal(data)
	if err != nil {
		log.Println("error ", err)
		return "", err
	}
	err = os.WriteFile(f.filename, out, 0644)
	if err != nil {
		log.Println("error ", err)
		return "", err
	}
	return "", nil
}

func (f *File) Find(_ *gin.Context, short string) (string, error) {
	fmt.Println("Findf")
	f.mx.Lock()
	defer f.mx.Unlock()
	fileData, err := os.ReadFile(f.filename)
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
	return "", nil
}

func (f *File) FindByUID(_ *gin.Context, uid string) ([]string, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	fileData, err := os.ReadFile(f.filename)
	if err != nil {
		log.Println("error ", err)
		return nil, err
	}
	parseData := []string{}
	json.Unmarshal(fileData, &parseData)
	answer := make([]string, 0, 4)
	for _, value := range parseData {
		var v = strings.Fields(value)
		if uid == v[2] {
			answer = append(answer, v[0]+" "+v[1])
		}
	}
	return answer, nil
}

func (f *File) Ping() (string, error) {
	return "Yes! Ping from File\n", nil
}

func (f *File) StoreBatch(_ *gin.Context, _ []map[string]string) error {
	return nil
}

func funnelf() {
	v := <-ch
	fmt.Println("funnelf v - ", v)
	Delf(v)
}

func Delf(p S) {
	fmt.Println("Delf - ", p)
}

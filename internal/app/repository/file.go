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

// !
type File struct {
	mx       sync.RWMutex
	filename string
}

// !
func NewFileStorage(f string) *File {
	p := &File{filename: f}
	go funnelf(p)
	return p
}

// !
func (f *File) Store(_ *gin.Context, original, uid string) (string, error) {
	short, err := GetRand()
	if err != nil {
		fmt.Println("error from GetRand")
	}
	data := make([]string, 0, 1)
	f.mx.Lock()
	defer f.mx.Unlock()
	_, err = os.Stat(f.filename)
	if err != nil {
		if os.IsNotExist(err) {
			os.OpenFile(f.filename, os.O_CREATE, 0777)
		}
		data = append(data, short+" "+original+" "+uid+" "+"true")
	} else {
		content, err := os.ReadFile(f.filename)
		if err != nil {
			return "", err
		}
		json.Unmarshal(content, &data)
		data = append(data, short+" "+original+" "+uid+" "+"true")
	}

	out, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(f.filename, out, 0644)
	if err != nil {
		return "", err
	}
	return short, nil
}

// !
func (f *File) Find(_ *gin.Context, short string) (string, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	fileData, err := os.ReadFile(f.filename)
	if err != nil {
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

// !
func (f *File) FindByUID(_ *gin.Context, uid string) ([]string, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	fileData, err := os.ReadFile(f.filename)
	if err != nil {
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

// !
func (f *File) Ping() (string, error) {
	return "Yes! Ping from File\n", nil
}

// !
func (f *File) StoreBatch(_ *gin.Context, batch []map[string]string) error {
	data := make([]string, 0, 1)
	for _, v := range batch {
		data = append(data, v["short_url"]+" "+v["original_url"]+" "+v["uid"]+" "+"true")
	}
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(f.filename, out, 0644)
	if err != nil {
		return err
	}
	return nil
}

func funnelf(f *File) {
	data := make([]string, 0, 1)
	for v := range ch {
		vch := v.Data
		uid := v.UID
		for _, s := range vch {
			content, err := os.ReadFile(f.filename)
			if err != nil {
				log.Println("error - ", err)
				return
			}
			json.Unmarshal(content, &data)
			for i, value := range data {
				var w = strings.Fields(value)
				if uid == w[2] && s == w[0] {
					data[i] = w[0] + " " + w[1] + " " + w[2] + " " + "false"
				}
			}
		}
		out, err := json.Marshal(data)
		if err != nil {
			log.Println("error ", err)
		}
		err = os.WriteFile(f.filename, out, 0644)
		if err != nil {
			log.Println("error ", err)
		}
	}
	Delf(f, data)
}

// !
func Delf(f *File, data []string) {

	for i, value := range data {
		var v = strings.Fields(value)
		if v[3] == "false" {
			data = append(data[:i], data[i+1:]...)
		}
	}
	out, err := json.Marshal(data)
	if err != nil {
		log.Println("error ", err)
	}
	err = os.WriteFile(f.filename, out, 0644)
	if err != nil {
		log.Println("error ", err)
	}
}

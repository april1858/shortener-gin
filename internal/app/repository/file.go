package repository

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func (r *Repository) FileStore(filename, short, original string) error {
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
		M = append(M, short+" "+original+" "+UID)
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
	return nil
}

func (r Repository) FileFind(filename, short string) (string, error) {
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
	return "", nil
}

func (r *Repository) FileFindByUID(filename string) ([]string, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error ", err)
		return nil, err
	}
	parseData := []string{}
	json.Unmarshal(fileData, &parseData)
	answer := make([]string, 0, 4)
	for _, value := range parseData {
		var v = strings.Fields(value)
		if UID == v[2] {
			answer = append(answer, v[0]+" "+v[1])
		}
	}
	return answer, nil
}

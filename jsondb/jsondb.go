package jsondb

import (
	"encoding/json"
	"os"
)

type JsonDB struct {
	fileName string
	content  map[string][]string
}

func NewJsonDB(fileName string) (*JsonDB, error) {
	f1, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &JsonDB{
				fileName: fileName,
				content:  make(map[string][]string),
			}, nil
		}
		return nil, err
	}
	defer f1.Close()
	simpleDB := make(map[string][]string)
	if err := json.NewDecoder(f1).Decode(&simpleDB); err != nil {
		simpleDB = make(map[string][]string)
	}
	return &JsonDB{
		fileName: fileName,
		content:  simpleDB,
	}, nil
}

func (j *JsonDB) Read(problem string) []string {
	dates, ok := j.content[problem]
	if !ok {
		return []string{}
	}
	return dates
}

func (j *JsonDB) Write(problem string, date string) error {
	f1, err := os.OpenFile(j.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f1.Close()
	dates, ok := j.content[problem]
	if !ok {
		dates = make([]string, 0, 1)
	}
	j.content[problem] = append(dates, date)
	encoder := json.NewEncoder(f1)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(j.content)
	return err
}

package jsonFileInteraction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Service interface {
	UnmarshalJSONFile(file string) error
	ReturnFileNamesArray() []FileName
}

type service struct {
	FileNames []FileName `json:"filenames"`
}

func NewJSONService() Service {
	json := &service{}
	json.UnmarshalJSONFile("config/fileNames.json")
	return json
}

type FileName struct {
	Name         string `json:"filename"`
	ReadableName string `json:"readablename"`
}

func (s *service) ReturnFileNamesArray() []FileName {
	return s.FileNames
}

func (s *service) UnmarshalJSONFile(file string) error {

	jsonFile, err := os.Open(file)

	if err != nil {
		return err
	}

	fmt.Println("Successfully Opened json file")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &s)

	return nil
}

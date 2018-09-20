package jsonFileInteraction

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
)

type Service interface {
	UnmarshalJSONFile(file string) error
	ReturnFileNamesArray() []FileName
	ReturnLocationsArray() []Location
	ReturnBackdatedFilesArray() []string
	ReturnAfterHoursFilesArray() []string
}

type service struct {
	FileNames       []FileName `json:"filenames"`
	Locations       []Location `json:"locations"`
	BackdatedFiles  []string   `json:"backdatedfiles"`
	AfterHoursFiles []string   `json:"afterhoursfiles"`
}

type FileName struct {
	Name         string `json:"filename"`
	ReadableName string `json:"readablename"`
}

type Location struct {
	Name      string   `json:"name"`
	MountPath string   `json:"mountpath"`
	Files     []string `json:"files"`
}

func NewJSONService() Service {
	json := &service{}
	json.UnmarshalJSONFile("/opt/app/fileNames.json")
	json.UnmarshalJSONFile("/opt/app/locations.json")
	return json
}

func (s *service) ReturnFileNamesArray() []FileName {
	return s.FileNames
}

func (s *service) ReturnLocationsArray() []Location {
	return s.Locations
}

func (s *service) ReturnBackdatedFilesArray() []string {
	return s.BackdatedFiles
}

func (s *service) ReturnAfterHoursFilesArray() []string {
	return s.AfterHoursFiles
}

func (s *service) UnmarshalJSONFile(file string) error {

	jsonFile, err := os.Open(file)

	if err != nil {
		return err
	}

	log.Printf("Successfully Opened %v\n", file)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &s)

	return nil
}

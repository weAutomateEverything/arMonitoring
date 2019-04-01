package jsonFileInteraction

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Service interface {
	UnmarshalJSONFile(file string) error
	ReturnFileNamesArray() []FileName
	ReturnLocationsArray() []Location
	ReturnBackdatedFilesArray() []string
	ReturnAfterHoursFilesArray() []string
	ReturnGenericFileNameArray() []string
	ReturnFileExpectedArrivalTimesMap() []ExpectedArrivalTimes
}

type service struct {
	FileNames                []FileName             `json:"filenames"`
	Locations                []Location             `json:"locations"`
	BackdatedFiles           []string               `json:"backdatedfiles"`
	AfterHoursFiles          []string               `json:"afterhoursfiles"`
	GenericFileNames         []string               `json:"genericnames"`
	FileExpectedArrivalTimes []ExpectedArrivalTimes `json:"expectedarrivaltimes"`
}

type FileName struct {
	Name         string `json:"filename"`
	ReadableName string `json:"readablename"`
}

type Location struct {
	Name      string   `json:"name"`
	TabNumber string   `json:"tabnumber"`
	MountPath string   `json:"mountpath"`
	Files     []string `json:"files"`
}

type ExpectedArrivalTimes struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

func NewJSONService() Service {
	json := &service{}
	json.UnmarshalJSONFile("/opt/app/fileNames.json")
	json.UnmarshalJSONFile("/opt/app/locations.json")
	json.UnmarshalJSONFile("/opt/app/genericNames.json")
	json.UnmarshalJSONFile("/opt/app/fileExpectedArrivalTime.json")
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

func (s *service) ReturnGenericFileNameArray() []string {
	return s.GenericFileNames
}

func (s *service) ReturnFileExpectedArrivalTimesMap() []ExpectedArrivalTimes {
	return s.FileExpectedArrivalTimes
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

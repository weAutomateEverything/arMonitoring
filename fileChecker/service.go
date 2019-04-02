package fileChecker

import (
	"fmt"
	"github.com/weAutomateEverything/fileMonitorService/jsonFileInteraction"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
	GetValues() map[string]string
	GetLocationName() string
	GetTabNumber() string
	Reset()
	ResetAfterHours()
	setValues(name, mountpath string, bdFiles, files []string, store Store)
	storeLocationStateRecent(name string, fileStatus map[string]string)
	setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error)
	getListOfFilesInPath(path string) ([]string, error)
	convertFileNamesToHumanReadableNames() map[string]string
}

type service struct {
	locationName    string
	tabNumber       string
	mountPath       string
	files           []string
	backDatedFiles  []string
	afterHoursFiles []string
	fileStatus      map[string]string
	store           Store
	json            jsonFileInteraction.Service
}

//Create New FileChecker instance
func NewFileChecker(json jsonFileInteraction.Service, tabNumber string, store Store, name, mountpath string, bdFiles []string, aHFiles []string, files ...string) Service {

	s := &service{
		locationName:    name,
		tabNumber:       tabNumber,
		mountPath:       mountpath,
		files:           files,
		backDatedFiles:  bdFiles,
		afterHoursFiles: aHFiles,
		fileStatus:      make(map[string]string),
		store:           store,
		json:            json,
	}

	storeContents, err := store.getLocationStateRecent(name)
	if err != nil {
		log.Println("Failed to access persistence layer with the following error: ", err)
	}

	if storeContents != nil {
		s.fileStatus = storeContents
	}

	go func() {
		s.setValues(s.locationName, s.mountPath, s.backDatedFiles, s.files, s.store)
	}()

	return s
}

func (s *service) GetValues() map[string]string {

	humanReadableFileStatusResponse := s.convertFileNamesToHumanReadableNames()

	return humanReadableFileStatusResponse
}

func (s *service) GetLocationName() string {

	return s.locationName
}

func (s *service) GetTabNumber() string {

	return s.tabNumber
}

func (s *service) Reset() {
	log.Printf("Resetting %s", s.locationName)
	for k := range s.fileStatus {
		if !isFileAfterHours(k, s.afterHoursFiles) {
			s.fileStatus[k] = "notreceived"
		}
	}
}

func (s *service) ResetAfterHours() {
	log.Printf("Resetting %s", s.locationName)
	for k := range s.fileStatus {
		if isFileAfterHours(k, s.afterHoursFiles) {
			s.fileStatus[k] = "notreceived"
		}
	}
}

func (s *service) setValues(name, mountpath string, bdFiles, files []string, store Store) {

	for true {
		log.Println(fmt.Sprintf("Now accessing %s share", name))

		notAccessableOrEmpty := isShareFolderEmpty(mountpath)

		for _, x := range files {
			//Execute the following if the share is inaccessible
			if notAccessableOrEmpty {
				for _, file := range files {
					s.fileStatus[file] = "unaccessable"
				}
				continue
			}
			//Retrieve the status for the current file in the loop
			value, err := s.setFileStatus(name, mountpath, x, bdFiles, store)
			if err != nil {
				log.Println(err)
			}
			//Reset inaccessible to notreceived once access has been re-established
			if _, ok := s.fileStatus[x]; ok {
				if s.fileStatus[x] == "unaccessable" {
					s.fileStatus[x] = "notreceived"
				}
				if s.fileStatus[x] == "late" || s.fileStatus[x] == "received" {
					continue
				}
			}
			//Set the value of the current file referenced in loop
			s.fileStatus[x] = value
		}
		//Store These latest notification stats in the Mongo DB
		s.storeLocationStateRecent(s.locationName, s.fileStatus)

		time.Sleep(1 * time.Minute)
	}
}

func (s *service) storeLocationStateRecent(name string, fileStatus map[string]string) {
	s.store.setLocationStateRecent(name, fileStatus)
}

func (s *service) setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error) {

	fileList, err := s.getListOfFilesInPath(dirPath)

	if err != nil {
		log.Println(fmt.Sprintf("Failed to access %s. Trying again in 2 seconds", dirPath))
	}

	yesterdayDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	currentDate := time.Now().Format("20060102")
	currentTime := time.Now().Format("15:04:05")

	convertedTime := convertTime("currentTime ", currentTime)

	for _, file := range fileList {

		backdated := isFileBackDated(file, bdFiles)

		expectedTime := s.expectedFileArivalTime(fileContains)
		contains := strings.Contains(file, fileContains)

		if backdated && contains {
			recent := strings.Contains(file, yesterdayDate)
			if recent && convertedTime.After(expectedTime) {
				return "late", nil
			}
			if recent && convertedTime.Before(expectedTime) {
				return "received", nil
			}
		} else if contains {
			recent := strings.Contains(file, currentDate)
			receivedYesterday := strings.Contains(file, yesterdayDate)

			if isFileAfterHours(file, s.afterHoursFiles) && receivedYesterday && s.isAfterHoursFileReceivedNextDay(fileContains) {
				return "late", nil
			}
			if recent && convertedTime.After(expectedTime) {
				return "late", nil
			}
			if recent && convertedTime.Before(expectedTime) {
				return "received", nil
			}
		}
	}
	return "notreceived", nil
}

func createHumanReadableResponseMap(notifications map[string]string) map[string]string {

	humanReadableFileStatusResponse := make(map[string]string)

	for k, v := range notifications {
		humanReadableFileStatusResponse[k] = v
	}

	return humanReadableFileStatusResponse
}

func (s *service) convertFileNamesToHumanReadableNames() map[string]string {

	genericFileNameArray := s.json.ReturnGenericFileNameArray()
	humanReadableFileNameList := s.json.ReturnFileNamesArray()

	humanReadableFileStatusResponse := createHumanReadableResponseMap(s.fileStatus)

	for k, v := range humanReadableFileStatusResponse {
		for _, fileName := range genericFileNameArray {
			if strings.Contains(k, fileName) {
				humanReadableFileStatusResponse[fileName] = v
				delete(humanReadableFileStatusResponse, k)
			}
		}
	}
	for _, fileName := range humanReadableFileNameList {
		if _, ok := s.fileStatus[fileName.Name]; ok {
			humanReadableFileStatusResponse[fileName.ReadableName] = humanReadableFileStatusResponse[fileName.Name]
			delete(humanReadableFileStatusResponse, fileName.Name)
		}
	}

	return humanReadableFileStatusResponse
}


func isShareFolderEmpty(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return true
	}
	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	defer dir.Close()

	return false
}

func isFileBackDated(file string, bdFiles []string) bool {

	var fileIsBackdated bool

	for _, bdfile := range bdFiles {
		if strings.Contains(file, bdfile) {
			fileIsBackdated = true
			continue
		}
	}
	return fileIsBackdated
}

func isFileAfterHours(file string, aHFiles []string) bool {

	var fileIsAfterHours bool

	for _, aHFile := range aHFiles {
		if strings.Contains(file, aHFile) {
			fileIsAfterHours = true
			continue
		}
	}
	return fileIsAfterHours
}

func convertTime(specifictime ,unconvertedTime string) time.Time {

	t, err := time.Parse("15:04:05", unconvertedTime)
	if err != nil {
		log.Printf("Failed to convert %v with the following error: %v",specifictime, err)
	}
	return t
}

func (s *service) isAfterHoursFileReceivedNextDay(file string) bool {
	currentTime := time.Now()
	if currentTime.Before(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 17, 0, 0, 0, currentTime.Location())) && currentTime.After(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 5, 0, 0, currentTime.Location())) && s.fileStatus[file] == "notreceived" {
		return true
	}
	return false
}

func (s *service) getListOfFilesInPath(path string) ([]string, error) {

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	list, _ := dir.Readdirnames(0)

	return list, nil
}

func (s *service) expectedFileArivalTime(file string) time.Time {

	var expectedTime string

	expectedTimes := s.json.ReturnFileExpectedArrivalTimesMap()

	for _, f := range expectedTimes {
		if f.Name == file {
			expectedTime = f.Time
		}
	}

	t := convertTime("expectedTime", expectedTime)

	return t
}

package fileAvailability

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"log"
	"github.com/jasonlvhit/gocron"
)

type Service interface {
	getFilesInPath(path string) ([]File, error)
	pathToMostRecentFile(dirPath, fileContains string) (string, time.Time, error)
	ConfirmUgandaFileAvailability()
}

type service struct {
}

func NewService() Service {
	return &service{}
}

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func (s *service) schedule() {
	confirmAvailability := gocron.NewScheduler()

	go func() {
		confirmAvailability.Every(1).Day().At("00:05").Do(s.ConfirmUgandaFileAvailability)
		<-confirmAvailability.Start()
	}()
}

func (s *service) getFilesInPath(path string) ([]File, error) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]File, len(files))

	for x, file := range files {
		result[x] = File{Name: file.Name(), LastModified: file.ModTime(), Path: path, Size: file.Size()}
	}
	return result, nil
}

func (s *service) pathToMostRecentFile(dirPath, fileContains string) (string, time.Time, error) {

	fileList, err := s.getFilesInPath(dirPath)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to access %v", dirPath))
	}

	currentDate := time.Now().Format("02/01/2006")

	for _, file := range fileList {
		cont := strings.Contains(file.Name, fileContains)

		daDate := file.LastModified.Format("02/01/2006")
		if daDate == currentDate && cont == true {
			return file.Name, file.LastModified, nil
		}
	}
	return "", time.Time{nil,nil, nil}, fmt.Errorf("%v file has not arrived yet", fileContains)
}

func (s *service) ConfirmUgandaFileAvailability() {
	fileName, fileModTime, err := s.pathToMostRecentFile("/mnt/uganda/",".txt")
	if err != nil{
		log.Println("Uganda file has not arrived yet")
	}
	log.Println(fmt.Sprintf("%v successfully received on %v", fileName, fileModTime))
}
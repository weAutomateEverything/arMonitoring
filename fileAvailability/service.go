package fileAvailability

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"github.com/matryer/try"
)

type Service interface {
	GetFilesInPath(path string) ([]File, error)
	pathToMostRecentFile(dirPath, fileContains string) (string, time.Time, error)
	ConfirmUgandaFileAvailability()
}

type service struct {
}

func NewService() Service {
	s := &service{}
	//s.ConfirmUgandaFileAvailability()
	return s
}

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func (s *service) schedule() {
	confirmUgandaAvailability := gocron.NewScheduler()

	go func() {
		confirmUgandaAvailability.Every(1).Day().At("00:05").Do(s.ConfirmUgandaFileAvailability)
		<-confirmUgandaAvailability.Start()
	}()
}

func (s *service) GetFilesInPath(path string) ([]File, error) {

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

	fileList, err := s.GetFilesInPath(dirPath)
	if err != nil || len(fileList) == 0{
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
	return "", time.Time{}, fmt.Errorf("%v file has not arrived yet", fileContains)
}

func (s *service) ConfirmFileAvailabilityMethod(path string) error{
	fileName, fileModTime, err := s.pathToMostRecentFile(path, ".TXT")
	if err != nil {
		return err
	} 
		log.Println(fmt.Sprintf("%v successfully received on %v", fileName, fileModTime))
	
}

func (s *service) ConfirmUgandaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 120
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/uganda")
		if err != nil {
			log.Println("Uganda file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}
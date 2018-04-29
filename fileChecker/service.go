package fileChecker

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
}

type service struct {
	mountPath  string
	fileStatus map[string]bool
}

func NewFileChecker(mountpath string, files ...string) *service {

	s := &service{
		mountPath:  mountpath,
		fileStatus: make(map[string]bool),
	}

	for _, x := range files {
		value := s.pathToMostRecentFile(mountpath, x)
		s.fileStatus[x] = value
	}

	//go func() {
	//	confirmAvailability := gocron.NewScheduler()
	//	confirmAvailability.Every(1).Minute().Do(s.ConfirmFileAvailabilityMethod)
	//}()
	
	return s
}

func (s *service) pathToMostRecentFile(dirPath, fileContains string) bool {

	fileList, err := s.GetFilesInPath(dirPath)
	if err != nil || len(fileList) == 0 {
		log.Println(fmt.Sprintf("Unable to access %v", dirPath))
	}

	currentDate := time.Now().Format("20060102")

	for _, file := range fileList {
		cont := strings.Contains(file, fileContains)
		recent := strings.Contains(file, currentDate)

		if recent == true && cont == true {
			return true
		}
	}
	return false
}

func (s *service) GetFilesInPath(path string) ([]string, error) {

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	list, _ := dir.Readdirnames(0)

	return list, nil
}

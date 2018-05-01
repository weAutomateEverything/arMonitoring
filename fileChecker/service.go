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
	mountPath            string
	fileStatus           map[string]bool
	fileStatusCollection map[string]map[string]bool
}

func NewFileChecker(mountpath string, files ...string) map[string]map[string]bool {

	s := &service{
		mountPath:            mountpath,
		fileStatus:           make(map[string]bool),
		fileStatusCollection: make(map[string]map[string]bool),
	}

	for _, x := range files {
		value := s.pathToMostRecentFile(mountpath, x)
		s.fileStatus[x] = value
	}
	s.fileStatusCollection[mountpath] = s.fileStatus

	return s.fileStatusCollection
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

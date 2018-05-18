package fileChecker

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"github.com/matryer/try"
)

type Service interface {
}

type service struct {
	locationName		 string
	mountPath            string
	fileStatus           map[string]bool
	fileStatusCollection map[string]map[string]bool
}

func NewFileChecker(name, mountpath string, files ...string) map[string]map[string]bool {

	s := &service{
		locationName:		  name,
		mountPath:            mountpath,
		fileStatus:           make(map[string]bool),
		fileStatusCollection: make(map[string]map[string]bool),
	}

	log.Println(fmt.Sprintf("Now accessing %s share", name))

	for _, x := range files {
		value := s.pathToMostRecentFile(mountpath, x)
		s.fileStatus[x] = value
	}
	s.fileStatusCollection[name] = s.fileStatus

	log.Println(fmt.Sprintf("Completed file confirmation process on %s share", name))

	return s.fileStatusCollection
}

func (s *service) pathToMostRecentFile(dirPath, fileContains string) bool {

	var fileList []string

	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 5
		var err error
		fileList, err = s.GetFilesInPath(dirPath)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to access %s. Trying again in 20 seconds", dirPath))
			time.Sleep(2 * time.Second)
		}
		return attempt < 5, err
	})
	if err != nil {
		log.Println(fmt.Sprintf("Unable to access %s", dirPath))
		return false
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

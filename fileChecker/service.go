package fileChecker

import (
	"fmt"
	"github.com/matryer/try"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
	GetValues() map[string]string
	GetLocationName() string
	Reset()
}

type service struct {
	locationName string
	mountPath    string
	files        []string
	fileStatus   map[string]string
	store        Store
}

func NewFileChecker(store Store, name, mountpath string, files ...string) Service {

	s := &service{
		locationName: name,
		mountPath:    mountpath,
		files:        files,
		fileStatus:   make(map[string]string),
		store:        store,
	}

	s.fileStatus = store.getLocationStateRecent(name)

	go func() {
		s.setValues(s.locationName, s.mountPath, s.files)
	}()

	return s
}

func (s *service) GetValues() map[string]string {
	return s.fileStatus
}

func (s *service) GetLocationName() string {
	return s.locationName
}

func (s *service) Reset() {
	log.Printf("Resetting %s", s.locationName)
	for k := range s.fileStatus {
		s.fileStatus[k] = "notreceived"
	}
}

func (s *service) setValues(name, mountpath string, files []string) {

	for true {
		log.Println(fmt.Sprintf("Now accessing %s share", name))

		for _, x := range files {
			value, err := s.setFileStatus(mountpath, x)
			if err != nil {
				for _, file := range files {
					s.fileStatus[file] = "unaccessable"
				}
			}
			if _, ok := s.fileStatus[x]; ok {
				if s.fileStatus[x] == "unaccessable" {
					s.fileStatus[x] = "notreceived"
				}
				if s.fileStatus[x] == "late" || s.fileStatus[x] == "received" {
					continue
				}
			}
			s.fileStatus[x] = value
		}

		s.storeLocationStateRecent(s.locationName, s.fileStatus)
		log.Println(fmt.Sprintf("Completed file confirmation process on %s share", name))
		time.Sleep(4 * time.Minute)
	}
}

func (s *service) storeLocationStateRecent(name string, fileStatus map[string]string) {
	s.store.addLocationStateRecent(name, fileStatus)
}

func (s *service) setFileStatus(dirPath, fileContains string) (string, error) {

	var fileList []string

	err := try.Do(func(attempt int) (bool, error) {

		var err error
		fileList, err = s.getListOfFilesInPath(dirPath)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to access %s. Trying again in 2 seconds", dirPath))
			time.Sleep(2 * time.Second)
		}
		return attempt < 5, err
	})
	if err != nil {
		log.Println(fmt.Sprintf("Unable to access %s, Please confirm share mount access", dirPath))
		return "", err
	}

	currentDate := time.Now().Format("20060102")
	currentTime := time.Now().Format("15:04:05")

	convertedTime := convertTime(currentTime)

	for _, file := range fileList {
		expectedTime := expectedFileArivalTime(fileContains)
		cont := strings.Contains(file, fileContains)
		recent := strings.Contains(file, currentDate)

		if recent == true && cont == true && convertedTime.After(expectedTime) {
			return "late", nil
		}
		if recent == true && cont == true && convertedTime.Before(expectedTime) {
			return "received", nil
		}
	}
	return "notreceived", nil
}

func convertTime(unconvertedTime string) time.Time {

	t, err := time.Parse("15:04:05", unconvertedTime)
	if err != nil {
		log.Printf("Failed to convert time with the following error: %v", err)
	}
	return t
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

func expectedFileArivalTime(file string) time.Time {

	expectedTimes := map[string]string{

		".001": "18:00:00",
		".002": "18:00:00",
		".003": "18:00:00",
		".004": "18:00:00",
		".005": "18:00:00",
		".006": "18:00:00",
		"DA":   "05:30:00",
		"DCI_OUTGOING_MONET_TRANS_REPORT": "05:30:00",
		"DCI_TRANS_INPUT_LIST_":           "05:30:00",
		"EP":                                "07:30:00",
		"GL":                                "01:30:00",
		"INT00001":                          "01:30:00",
		"INT00003":                          "01:30:00",
		"INT00007":                          "01:30:00",
		"MASTERCARD_ACKNOWLEDGEMENT_REPORT": "05:30:00",
		"MC_INCOMING_FILES_SUMMARY_REPORT":  "05:30:00",
		"MC_INCOMING_MONET_TRANS_REPORT":    "05:30:00",
		"MC_OUTGOING_FILES_SUMMARY_REPORT":  "05:30:00",
		"MC_OUTGOING_MONET_TRANS_REPORT":    "05:30:00",
		"MERCH_REJ_TRANS":                   "05:30:00",
		"MS":                                "05:30:00",
		"MUL":                               "01:30:00",
		"RECON_REPORT":                      "05:30:00",
		"SE":                                "01:30:00",
		"SPTLSB_NA_":                        "21:00:00",
		"SR00001":                           "01:30:00",
		"TRANS_INPUT_LIST_":                 "05:30:00",
		"TXN":                               "01:30:00",
		"VISA_INCOMING_FILES_SUMMARY_REPORT": "05:30:00",
		"VISA_INCOMING_MONET_TRANS_REPORT":   "05:30:00",
		"VISA_OUTGOING_FILES_SUMMARY_REPORT": "05:30:00",
		"VISA_OUTGOING_MONET_TRANS_REPORT":   "05:30:00",
		"VOUT":  "07:30:00",
		"VTRAN": "07:30:00",
		"EP747": "07:30:00",
		"CGNI":  "07:30:00",
	}

	expectedTime := expectedTimes[file]

	t := convertTime(expectedTime)

	return t
}

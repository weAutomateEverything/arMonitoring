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
}

type service struct {
	locationName string
	mountPath    string
	fileStatus   map[string]string
}

func NewFileChecker(name, mountpath string, files ...string) (string, map[string]string) {

	s := &service{
		locationName: name,
		mountPath:    mountpath,
		fileStatus:   make(map[string]string),
	}

	log.Println(fmt.Sprintf("Now accessing %s share", name))

	for _, x := range files {
		value, err := s.pathToMostRecentFile(mountpath, x)
		if err != nil {
			for _, file := range files {
				s.fileStatus[file] = "unaccessable"
			}
			return s.locationName, s.fileStatus
		}
		s.fileStatus[x] = value
	}

	log.Println(fmt.Sprintf("Completed file confirmation process on %s share", name))

	return s.locationName, s.fileStatus
}

func (s *service) pathToMostRecentFile(dirPath, fileContains string) (string, error) {

	var fileList []string

	err := try.Do(func(attempt int) (bool, error) {

		var err error
		fileList, err = s.GetFilesInPath(dirPath)
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

func expectedFileArivalTime(file string) time.Time {

	expectedTimes := map[string]string{

		".001": "01:30:00",
		".002": "01:30:00",
		".003": "01:30:00",
		".004": "01:30:00",
		".005": "01:30:00",
		".006": "01:30:00",
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

func convertTime(unconvertedTime string) time.Time {

	t, err := time.Parse("15:04:05", unconvertedTime)
	if err != nil {
		log.Printf("Failed to convert time with the following error: %v", err)
	}
	return t
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

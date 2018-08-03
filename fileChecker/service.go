package fileChecker

import (
	"fmt"
	"github.com/matryer/try"
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
	Reset()
	ResetAfterHours()
	setValues(name, mountpath string, bdFiles []string, files []string, store Store)
	storeLocationStateRecent(name string, fileStatus map[string]string)
	setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error)
	getListOfFilesInPath(path string) ([]string, error)
	convertFileNamesToHumanReadableNames()
}

type service struct {
	locationName    string
	mountPath       string
	files           []string
	backDatedFiles  []string
	afterHoursFiles []string
	fileStatus      map[string]string
	store           Store
	json            jsonFileInteraction.Service
}

func NewFileChecker(json jsonFileInteraction.Service, store Store, name, mountpath string, bdFiles []string, aHFiles []string, files ...string) Service {

	s := &service{
		locationName:    name,
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
		log.Println("Failed to access persistance layer with the following error: ", err)
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

	//Convert names to human readable
	s.convertFileNamesToHumanReadableNames()

	return s.fileStatus
}

func (s *service) GetLocationName() string {
	return s.locationName
}

func (s *service) Reset() {
	log.Printf("Resetting %s", s.locationName)
	for k := range s.fileStatus {
		if !checkIfFileIsAfterHours(k, s.afterHoursFiles) {
			s.fileStatus[k] = "notreceived"
		}
	}
}
func (s *service) ResetAfterHours() {
	log.Printf("Resetting %s", s.locationName)
	for k := range s.fileStatus {
		if checkIfFileIsAfterHours(k, s.afterHoursFiles) {
			s.fileStatus[k] = "notreceived"
		}
	}
}

func (s *service) setValues(name, mountpath string, bdFiles []string, files []string, store Store) {

	for true {
		log.Println(fmt.Sprintf("Now accessing %s share", name))

		notAccessableOrEmpty := isShareFolderEmpty(mountpath)

		for _, x := range files {
			if notAccessableOrEmpty {
				for _, file := range files {
					s.fileStatus[file] = "unaccessable"
				}
				continue
			}
			value, err := s.setFileStatus(name, mountpath, x, bdFiles, store)
			if err != nil {
				log.Println(err)
			}
			// Test for existance of key in map
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
	s.store.setLocationStateRecent(name, fileStatus)
}

func (s *service) setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error) {

	var fileList []string

	//Attempt connection
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

	yesterdayDate := time.Now().AddDate(0, 0, -1).Format("20060102")
	currentDate := time.Now().Format("20060102")
	currentTime := time.Now().Format("15:04:05")

	convertedTime := convertTime(currentTime)

	for _, file := range fileList {

		backdated := checkIfFileIsBackDated(file, bdFiles)

		expectedTime := expectedFileArivalTime(fileContains)
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

func (s *service) convertFileNamesToHumanReadableNames() {

	humanReadableFileNameList := s.json.ReturnFileNamesArray()

	for _, fileName := range humanReadableFileNameList {
		if _, ok := s.fileStatus[fileName.Name]; ok {
			s.fileStatus[fileName.ReadableName] = s.fileStatus[fileName.Name]
			delete(s.fileStatus, fileName.Name)
		}
	}
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

func checkIfFileIsBackDated(file string, bdFiles []string) bool {

	var fileIsBackdated bool

	for _, bdfile := range bdFiles {
		if strings.Contains(file, bdfile) {
			fileIsBackdated = true
			continue
		}
	}
	return fileIsBackdated
}

func checkIfFileIsAfterHours(file string, aHFiles []string) bool {

	var fileIsAfterHours bool

	for _, aHFile := range aHFiles {
		if strings.Contains(file, aHFile) {
			fileIsAfterHours = true
			continue
		}
	}
	return fileIsAfterHours
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
		"SPTLSB":                            "21:00:00",
		"SR00001":                           "01:30:00",
		"TRANS_INPUT_LIST_":                 "05:30:00",
		"TXN":                               "01:30:00",
		"VISA_INCOMING_FILES_SUMMARY_REPORT": "05:30:00",
		"VISA_INCOMING_MONET_TRANS_REPORT":   "05:30:00",
		"VISA_OUTGOING_FILES_SUMMARY_REPORT": "05:30:00",
		"VISA_OUTGOING_MONET_TRANS_REPORT":   "05:30:00",
		"VOUT":    "07:30:00",
		"VTRAN":   "07:30:00",
		"EP747":   "07:30:00",
		"CGNI":    "07:30:00",
		"GL149":   "01:30:00",
		"GL150":   "01:30:00",
		"DA149":   "01:30:00",
		"DA150":   "01:30:00",
		"SE149":   "01:30:00",
		"SE150":   "01:30:00",
		"MS149":   "01:30:00",
		"MS150":   "01:30:00",
		"149_TXN": "01:30:00",
		"150_TXN": "01:30:00",
	}

	expectedTime := expectedTimes[file]

	t := convertTime(expectedTime)

	return t
}

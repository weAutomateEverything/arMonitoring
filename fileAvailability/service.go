package fileAvailability

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/matryer/try"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Service interface {
	GetFilesInPath(path string) ([]File, error)
	pathToMostRecentFile(dirPath, fileContains string) (bool, string, time.Time, error)
	ConfirmUgandaFileAvailability()
	CreateJSONResponse() map[string]bool
}

type service struct {
}

var (
	ZimbabweStatus   bool
	BotswanaStatus   bool
	KenyaStatus      bool
	MalawiStatus     bool
	NamibiaStatus    bool
	GhanaStatus      bool
	UgandaStatus     bool
	UgandaDRStatus   bool
	ZambiaStatus     bool
	ZambiaDRStatus   bool
	ZambiaProdStatus bool
)

//type FileStatus struct {
//	location     string
//	path         string
//	received     bool
//	timeReceived time.Time
//}

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func NewService() Service {
	s := &service{}
	s.ConfirmUgandaFileAvailability()
	s.ConfirmBotswanaFileAvailability()
	s.ConfirmGhanaFileAvailability()
	s.ConfirmKenyaFileAvailability()
	s.ConfirmMalawiFileAvailability()
	s.ConfirmNamibiaFileAvailability()
	s.ConfirmUgandaDRFileAvailability()
	s.ConfirmZambiaFileAvailability()
	s.ConfirmZambiaDRFileAvailability()
	s.ConfirmZambiaProdFileAvailability()
	s.ConfirmZimbabweFileAvailability()
	
	return s
}

func (s *service) schedule() {
	confirmUgandaAvailability := gocron.NewScheduler()

	go func() {
		confirmUgandaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmUgandaFileAvailability)
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

func (s *service) pathToMostRecentFile(dirPath, fileContains string) (bool, string, time.Time, error) {

	fileList, err := s.GetFilesInPath(dirPath)
	if err != nil || len(fileList) == 0 {
		log.Println(fmt.Sprintf("Unable to access %v", dirPath))
	}

	currentDate := time.Now().Format("02/01/2006")

	for _, file := range fileList {
		cont := strings.Contains(file.Name, fileContains)

		daDate := file.LastModified.Format("02/01/2006")
		if daDate == currentDate && cont == true {
			return true, file.Name, file.LastModified, nil
		}
	}
	return false, "", time.Time{}, fmt.Errorf("%v file has not arrived yet", fileContains)
}

func (s *service) ConfirmFileAvailabilityMethod(path string) error {
	fileReceived, _, _, err := s.pathToMostRecentFile(path, ".txt")

	if err != nil {
		return err
	}

	switch path {
	case "/mnt/zimbabwe":
		ZimbabweStatus = fileReceived
	case "/mnt/botswana":
		BotswanaStatus = fileReceived
	case "/mnt/ghana":
		GhanaStatus = fileReceived
	case "/mnt/kenya":
		KenyaStatus = fileReceived
	case "/mnt/malawi":
		MalawiStatus = fileReceived
	case "/mnt/namibia":
		NamibiaStatus = fileReceived
	case "/mnt/uganda":
		UgandaStatus = fileReceived
	case "/mnt/ugandadr":
		UgandaDRStatus = fileReceived
	case "/mnt/zambia":
		ZambiaStatus = fileReceived
	case "/mnt/zambiadr":
		ZambiaDRStatus = fileReceived
	case "/mnt/zambiaprod":
		ZambiaProdStatus = fileReceived

	}
	return nil
}

func (s *service) CreateJSONResponse() map[string]bool {

	resp := map[string]bool{
		"ZimbabweStatus":   ZimbabweStatus,
		"BotswanaStatus":   BotswanaStatus,
		"KenyaStatus":      KenyaStatus,
		"MalawiStatus":     MalawiStatus,
		"NamibiaStatus":    NamibiaStatus,
		"GhanaStatus":      GhanaStatus,
		"UgandaStatus":     UgandaStatus,
		"UgandaDRStatus":   UgandaDRStatus,
		"ZambiaStatus":     ZambiaStatus,
		"ZambiaDRStatus":   ZambiaDRStatus,
		"ZambiaProdStatus": ZambiaProdStatus,
	}

	return resp
}

func (s *service) ConfirmZimbabweFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/zimbabwe")
		if err != nil {
			log.Println("Zimbabwe file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmBotswanaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/botswana")
		if err != nil {
			log.Println("Botswana file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmGhanaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/ghana")
		if err != nil {
			log.Println("Ghana file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmKenyaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/kenya")
		if err != nil {
			log.Println("Kenya file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmMalawiFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/malawi")
		if err != nil {
			log.Println("Malawi file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmNamibiaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/namibia")
		if err != nil {
			log.Println("Namibia file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmUgandaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
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

func (s *service) ConfirmUgandaDRFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/ugandadr")
		if err != nil {
			log.Println("UgandaDR file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmZambiaFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/zambia")
		if err != nil {
			log.Println("Zambia file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmZambiaDRFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/zambiadr")
		if err != nil {
			log.Println("ZambiaDR file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *service) ConfirmZambiaProdFileAvailability() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 240
		var err error
		err = s.ConfirmFileAvailabilityMethod("/mnt/zambiaprod")
		if err != nil {
			log.Println("ZambiaProd file not yet detected. Next attempt in 2 minutes...")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
	}
}

package fileAvailability

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/matryer/try"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
	GetFilesInPath(path string) ([]string, error)
	pathToMostRecentFile(dirPath, fileContains string) (bool, string)
	CreateJSONResponse() LocationStatus
	fileCount(dirPath, fileContains string) int
	ConfirmZimbabweFileAvailability()
}

type service struct {
}

type LocationStatus struct {
	ZimbabweStatus   FilesReceived
	BotswanaStatus   FilesReceived
	KenyaStatus      FilesReceived
	MalawiStatus     FilesReceived
	NamibiaStatus    FilesReceived
	GhanaStatus      FilesReceived
	UgandaStatus     FilesReceived
	UgandaDRStatus   FilesReceived
	ZambiaStatus     FilesReceived
	ZambiaDRStatus   FilesReceived
	ZambiaProdStatus FilesReceived
}

type FilesReceived struct {
	SE147    bool
	GL147    bool
	MUL00004 bool
	SR00001  bool
	TXN      bool
	DA147    bool
	MS147    bool
	EP747    bool
	VTRAN147 bool
	VOUT147  bool
	SPTLSB   bool
	CGNI     bool
	INT00001 bool
	INT00003 bool
	INT00007 bool
	PDF      int
	TT140    int
}

var locationStatus LocationStatus

func NewService() Service {
	s := &service{}
	//s.ConfirmBotswanaFileAvailability()
	//s.ConfirmGhanaFileAvailability()
	//s.ConfirmKenyaFileAvailability()
	//s.ConfirmMalawiFileAvailability()
	//s.ConfirmNamibiaFileAvailability()
	go s.ConfirmUgandaFileAvailability()
	//s.ConfirmUgandaDRFileAvailability()
	//s.ConfirmZambiaFileAvailability()
	//s.ConfirmZambiaDRFileAvailability()
	//s.ConfirmZambiaProdFileAvailability()
	go s.ConfirmZimbabweFileAvailability()

	return s
}

func (s *service) schedule() {
	confirmUgandaAvailability := gocron.NewScheduler()
	confirmBotswanaAvailability := gocron.NewScheduler()
	confirmMalawiAvailability := gocron.NewScheduler()
	confirmGhanaAvailability := gocron.NewScheduler()
	confirmKenyaAvailability := gocron.NewScheduler()
	confirmNamibiaAvailability := gocron.NewScheduler()
	confirmUgandaDRAvailability := gocron.NewScheduler()
	confirmZambiaAvailability := gocron.NewScheduler()
	confirmZambiaDRAvailability := gocron.NewScheduler()
	confirmZambiaProdAvailability := gocron.NewScheduler()
	confirmZimbabweAvailability := gocron.NewScheduler()
	resetStatus := gocron.NewScheduler()

	go func() {
		confirmUgandaAvailability.Every(1).Minute().Do(s.ConfirmUgandaFileAvailability)
		<-confirmUgandaAvailability.Start()
	}()
	go func() {
		confirmBotswanaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmBotswanaFileAvailability)
		<-confirmBotswanaAvailability.Start()
	}()
	go func() {
		confirmMalawiAvailability.Every(1).Day().At("00:00").Do(s.ConfirmMalawiFileAvailability)
		<-confirmMalawiAvailability.Start()
	}()
	go func() {
		confirmGhanaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmGhanaFileAvailability)
		<-confirmGhanaAvailability.Start()
	}()
	go func() {
		confirmKenyaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmKenyaFileAvailability)
		<-confirmKenyaAvailability.Start()
	}()
	go func() {
		confirmNamibiaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmNamibiaFileAvailability)
		<-confirmNamibiaAvailability.Start()
	}()
	go func() {
		confirmUgandaDRAvailability.Every(1).Day().At("00:00").Do(s.ConfirmUgandaDRFileAvailability)
		<-confirmUgandaDRAvailability.Start()
	}()
	go func() {
		confirmZambiaAvailability.Every(1).Day().At("00:00").Do(s.ConfirmZambiaFileAvailability)
		<-confirmZambiaAvailability.Start()
	}()
	go func() {
		confirmZambiaDRAvailability.Every(1).Day().At("00:00").Do(s.ConfirmZambiaDRFileAvailability)
		<-confirmZambiaDRAvailability.Start()
	}()
	go func() {
		confirmZambiaProdAvailability.Every(1).Day().At("00:00").Do(s.ConfirmZambiaProdFileAvailability)
		<-confirmZambiaProdAvailability.Start()
	}()
	go func() {
		confirmZimbabweAvailability.Every(1).Minute().Do(s.ConfirmZimbabweFileAvailability)
		<-confirmZimbabweAvailability.Start()
	}()
	go func() {
		resetStatus.Every(1).Day().At("00:00").Do(s.resetStatus)
		<-resetStatus.Start()
	}()
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

func (s *service) pathToMostRecentFile(dirPath, fileContains string) (bool, string) {

	fileList, err := s.GetFilesInPath(dirPath)
	if err != nil || len(fileList) == 0 {
		log.Println(fmt.Sprintf("Unable to access %v", dirPath))
	}

	currentDate := time.Now().Format("20060102")

	for _, file := range fileList {
		cont := strings.Contains(file, fileContains)
		recent := strings.Contains(file, currentDate)

		if recent == true && cont == true {
			return true, file
		}
	}
	return false, ""
}

func (s *service) fileCount(dirPath, fileContains string) int {

	fileCount := 0

	fileList, err := s.GetFilesInPath(dirPath)
	if err != nil || len(fileList) == 0 {
		log.Println(fmt.Sprintf("Unable to access %v", dirPath))
	}

	currentDate := time.Now().Format("20060102")

	for _, file := range fileList {
		cont := strings.Contains(strings.ToUpper(file), fileContains)
		recent := strings.Contains(file, currentDate)

		if recent == true && cont == true {
			fileCount++
		}
	}
	return fileCount
}

func (s *service) ConfirmFileAvailabilityMethod(path string) error {

	SE147, _ := s.pathToMostRecentFile(path, "SE")
	GL147, _ := s.pathToMostRecentFile(path, "GL")
	MUL00004, _ := s.pathToMostRecentFile(path, "MUL")
	SR00001, _ := s.pathToMostRecentFile(path, "SR")
	TXN, _ := s.pathToMostRecentFile(path, "TXN")
	DA147, _ := s.pathToMostRecentFile(path, "DA")
	MS147, _ := s.pathToMostRecentFile(path, "MS")
	EP747, _ := s.pathToMostRecentFile(path, "EP747")
	VTRAN147, _ := s.pathToMostRecentFile(path, "VTRAN")
	VOUT147, _ := s.pathToMostRecentFile(path, "VOUT")
	SPTLSB, _ := s.pathToMostRecentFile(path, "SPTLSB")
	CGNI, _ := s.pathToMostRecentFile(path, "CGNI")
	INT00001, _ := s.pathToMostRecentFile(path, "INT00001")
	INT00003, _ := s.pathToMostRecentFile(path, "INT00003")
	INT00007, _ := s.pathToMostRecentFile(path, "INT00007")

	PDF := s.fileCount(path, "PDF")
	TT140 := s.fileCount(path, "TT140")

	switch path {
	case "/mnt/zimbabwe":
		locationStatus.ZimbabweStatus.SE147 = SE147
		locationStatus.ZimbabweStatus.GL147 = GL147
		locationStatus.ZimbabweStatus.TXN = TXN
		locationStatus.ZimbabweStatus.VTRAN147 = VTRAN147
		locationStatus.ZimbabweStatus.VOUT147 = VOUT147
		locationStatus.ZimbabweStatus.MS147 = MS147
		locationStatus.ZimbabweStatus.DA147 = DA147
		locationStatus.ZimbabweStatus.EP747 = EP747
		locationStatus.ZimbabweStatus.PDF = PDF
		locationStatus.ZimbabweStatus.TT140 = TT140
	case "/mnt/botswana":
		locationStatus.BotswanaStatus.SE147 = SE147
		locationStatus.BotswanaStatus.GL147 = GL147
		locationStatus.BotswanaStatus.TXN = TXN
		locationStatus.BotswanaStatus.MUL00004 = MUL00004
		locationStatus.BotswanaStatus.VTRAN147 = VTRAN147
		locationStatus.BotswanaStatus.VOUT147 = VOUT147
		locationStatus.BotswanaStatus.MS147 = MS147
		locationStatus.BotswanaStatus.DA147 = DA147
		locationStatus.BotswanaStatus.EP747 = EP747
	case "/mnt/ghana":
		locationStatus.GhanaStatus.SE147 = SE147
		locationStatus.GhanaStatus.GL147 = GL147
		locationStatus.GhanaStatus.TXN = TXN
		locationStatus.GhanaStatus.MUL00004 = MUL00004
		locationStatus.GhanaStatus.VTRAN147 = VTRAN147
		locationStatus.GhanaStatus.VOUT147 = VOUT147
		locationStatus.GhanaStatus.MS147 = MS147
		locationStatus.GhanaStatus.DA147 = DA147
		locationStatus.GhanaStatus.EP747 = EP747
	case "/mnt/kenya":
		locationStatus.KenyaStatus.SE147 = SE147
		locationStatus.KenyaStatus.GL147 = GL147
		locationStatus.KenyaStatus.TXN = TXN
		locationStatus.KenyaStatus.MUL00004 = MUL00004
		locationStatus.KenyaStatus.VTRAN147 = VTRAN147
		locationStatus.KenyaStatus.VOUT147 = VOUT147
		locationStatus.KenyaStatus.MS147 = MS147
		locationStatus.KenyaStatus.DA147 = DA147
		locationStatus.KenyaStatus.EP747 = EP747
	case "/mnt/malawi":
		locationStatus.MalawiStatus.SE147 = SE147
		locationStatus.MalawiStatus.GL147 = GL147
		locationStatus.MalawiStatus.TXN = TXN
		locationStatus.MalawiStatus.MUL00004 = MUL00004
		locationStatus.MalawiStatus.VTRAN147 = VTRAN147
		locationStatus.MalawiStatus.VOUT147 = VOUT147
		locationStatus.MalawiStatus.MS147 = MS147
		locationStatus.MalawiStatus.DA147 = DA147
		locationStatus.MalawiStatus.EP747 = EP747
	case "/mnt/namibia":
		locationStatus.NamibiaStatus.SE147 = SE147
		locationStatus.NamibiaStatus.GL147 = GL147
		locationStatus.NamibiaStatus.TXN = TXN
		locationStatus.NamibiaStatus.MUL00004 = MUL00004
		locationStatus.NamibiaStatus.VTRAN147 = VTRAN147
		locationStatus.NamibiaStatus.VOUT147 = VOUT147
		locationStatus.NamibiaStatus.MS147 = MS147
		locationStatus.NamibiaStatus.DA147 = DA147
		locationStatus.NamibiaStatus.INT00001 = INT00001
		locationStatus.NamibiaStatus.INT00003 = INT00003
		locationStatus.NamibiaStatus.INT00007 = INT00007
		locationStatus.NamibiaStatus.SR00001 = SR00001
		locationStatus.NamibiaStatus.EP747 = EP747
		locationStatus.NamibiaStatus.SPTLSB = SPTLSB
		locationStatus.NamibiaStatus.CGNI = CGNI
	case "/mnt/uganda":
		locationStatus.UgandaStatus.SE147 = SE147
		locationStatus.UgandaStatus.GL147 = GL147
		locationStatus.UgandaStatus.TXN = TXN
		locationStatus.UgandaStatus.VTRAN147 = VTRAN147
		locationStatus.UgandaStatus.VOUT147 = VOUT147
		locationStatus.UgandaStatus.MS147 = MS147
		locationStatus.UgandaStatus.DA147 = DA147
		locationStatus.UgandaStatus.EP747 = EP747
		locationStatus.UgandaStatus.PDF = PDF
		locationStatus.UgandaStatus.TT140 = TT140
	case "/mnt/ugandadr":
		locationStatus.UgandaDRStatus.SE147 = SE147
		locationStatus.UgandaDRStatus.GL147 = GL147
		locationStatus.UgandaDRStatus.TXN = TXN
		locationStatus.UgandaDRStatus.VTRAN147 = VTRAN147
		locationStatus.UgandaDRStatus.VOUT147 = VOUT147
		locationStatus.UgandaDRStatus.MS147 = MS147
		locationStatus.UgandaDRStatus.DA147 = DA147
		locationStatus.UgandaDRStatus.EP747 = EP747
	case "/mnt/zambia":
		locationStatus.ZambiaStatus.SE147 = SE147
		locationStatus.ZambiaStatus.GL147 = GL147
		locationStatus.ZambiaStatus.TXN = TXN
		locationStatus.ZambiaStatus.VTRAN147 = VTRAN147
		locationStatus.ZambiaStatus.VOUT147 = VOUT147
		locationStatus.ZambiaStatus.MS147 = MS147
		locationStatus.ZambiaStatus.DA147 = DA147
		locationStatus.ZambiaStatus.EP747 = EP747
	case "/mnt/zambiadr":
		locationStatus.ZambiaDRStatus.SE147 = SE147
		locationStatus.ZambiaDRStatus.GL147 = GL147
		locationStatus.ZambiaDRStatus.TXN = TXN
		locationStatus.ZambiaDRStatus.VTRAN147 = VTRAN147
		locationStatus.ZambiaDRStatus.VOUT147 = VOUT147
		locationStatus.ZambiaDRStatus.MS147 = MS147
		locationStatus.ZambiaDRStatus.DA147 = DA147
		locationStatus.ZambiaDRStatus.EP747 = EP747
	case "/mnt/zambiaprod":
		locationStatus.ZambiaProdStatus.SE147 = SE147
		locationStatus.ZambiaProdStatus.GL147 = GL147
		locationStatus.ZambiaProdStatus.TXN = TXN
		locationStatus.ZambiaProdStatus.VTRAN147 = VTRAN147
		locationStatus.ZambiaProdStatus.VOUT147 = VOUT147
		locationStatus.ZambiaProdStatus.MS147 = MS147
		locationStatus.ZambiaProdStatus.DA147 = DA147
		locationStatus.ZambiaProdStatus.EP747 = EP747

	}
	return nil
}

func (s *service) CreateJSONResponse() LocationStatus {

	return locationStatus
}

func (s *service) resetStatus() {
	locationStatus = LocationStatus{}
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

package fileAvailability

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
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
	SE       bool
	GL       bool
	MUL      bool
	SR       bool
	TXN      bool
	DA       bool
	MS       bool
	EP    bool
	VTRAN    bool
	VOUT     bool
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
	//go s.ConfirmBotswanaFileAvailability()
	//go s.ConfirmGhanaFileAvailability()
	//go s.ConfirmKenyaFileAvailability()
	go s.ConfirmMalawiFileAvailability()
	//go s.ConfirmNamibiaFileAvailability()
	go s.ConfirmUgandaFileAvailability()
	//go s.ConfirmUgandaDRFileAvailability()
	go s.ConfirmZambiaFileAvailability()
	//go s.ConfirmZambiaDRFileAvailability()
	//go s.ConfirmZambiaProdFileAvailability()
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

	go func() {
		confirmUgandaAvailability.Every(1).Minute().Do(s.ConfirmUgandaFileAvailability)
		<-confirmUgandaAvailability.Start()
	}()
	go func() {
		confirmBotswanaAvailability.Every(1).Minute().Do(s.ConfirmBotswanaFileAvailability)
		<-confirmBotswanaAvailability.Start()
	}()
	go func() {
		confirmMalawiAvailability.Every(1).Minute().Do(s.ConfirmMalawiFileAvailability)
		<-confirmMalawiAvailability.Start()
	}()
	go func() {
		confirmGhanaAvailability.Every(1).Minute().Do(s.ConfirmGhanaFileAvailability)
		<-confirmGhanaAvailability.Start()
	}()
	go func() {
		confirmKenyaAvailability.Every(1).Minute().Do(s.ConfirmKenyaFileAvailability)
		<-confirmKenyaAvailability.Start()
	}()
	go func() {
		confirmNamibiaAvailability.Every(1).Minute().Do(s.ConfirmNamibiaFileAvailability)
		<-confirmNamibiaAvailability.Start()
	}()
	go func() {
		confirmUgandaDRAvailability.Every(1).Minute().Do(s.ConfirmUgandaDRFileAvailability)
		<-confirmUgandaDRAvailability.Start()
	}()
	go func() {
		confirmZambiaAvailability.Every(1).Minute().Do(s.ConfirmZambiaFileAvailability)
		<-confirmZambiaAvailability.Start()
	}()
	go func() {
		confirmZambiaDRAvailability.Every(1).Minute().Do(s.ConfirmZambiaDRFileAvailability)
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

	SE, _ := s.pathToMostRecentFile(path, "SE")
	GL, _ := s.pathToMostRecentFile(path, "GL")
	MUL, _ := s.pathToMostRecentFile(path, "MUL")
	SR, _ := s.pathToMostRecentFile(path, "SR")
	TXN, _ := s.pathToMostRecentFile(path, "TXN")
	DA, _ := s.pathToMostRecentFile(path, "DA")
	MS, _ := s.pathToMostRecentFile(path, "MS")
	EP, _ := s.pathToMostRecentFile(path, "EP")
	VTRAN, _ := s.pathToMostRecentFile(path, "VTRAN")
	VOUT, _ := s.pathToMostRecentFile(path, "VOUT")
	SPTLSB, _ := s.pathToMostRecentFile(path, "SPTLSB")
	CGNI, _ := s.pathToMostRecentFile(path, "CGNI")
	INT00001, _ := s.pathToMostRecentFile(path, "INT00001")
	INT00003, _ := s.pathToMostRecentFile(path, "INT00003")
	INT00007, _ := s.pathToMostRecentFile(path, "INT00007")

	PDF := s.fileCount(path, "PDF")
	TT140 := s.fileCount(path, "TT140")

	switch path {
	case "/mnt/zimbabwe":
		locationStatus.ZimbabweStatus.SE = SE
		locationStatus.ZimbabweStatus.GL = GL
		locationStatus.ZimbabweStatus.TXN = TXN
		locationStatus.ZimbabweStatus.MUL = MUL
		locationStatus.ZimbabweStatus.VTRAN = VTRAN
		locationStatus.ZimbabweStatus.VOUT = VOUT
		locationStatus.ZimbabweStatus.MS = MS
		locationStatus.ZimbabweStatus.DA = DA
		locationStatus.ZimbabweStatus.INT00001 = INT00001
		locationStatus.ZimbabweStatus.INT00003 = INT00003
		locationStatus.ZimbabweStatus.INT00007 = INT00007
		locationStatus.ZimbabweStatus.SR = SR
		locationStatus.ZimbabweStatus.EP = EP
		locationStatus.ZimbabweStatus.SPTLSB = SPTLSB
		locationStatus.ZimbabweStatus.CGNI = CGNI
		locationStatus.ZimbabweStatus.PDF = PDF
		locationStatus.ZimbabweStatus.TT140 = TT140
	case "/mnt/botswana":
		locationStatus.BotswanaStatus.SE = SE
		locationStatus.BotswanaStatus.GL = GL
		locationStatus.BotswanaStatus.TXN = TXN
		locationStatus.BotswanaStatus.MUL = MUL
		locationStatus.BotswanaStatus.VTRAN = VTRAN
		locationStatus.BotswanaStatus.VOUT = VOUT
		locationStatus.BotswanaStatus.MS = MS
		locationStatus.BotswanaStatus.DA = DA
		locationStatus.BotswanaStatus.INT00001 = INT00001
		locationStatus.BotswanaStatus.INT00003 = INT00003
		locationStatus.BotswanaStatus.INT00007 = INT00007
		locationStatus.BotswanaStatus.SR = SR
		locationStatus.BotswanaStatus.EP = EP
		locationStatus.BotswanaStatus.SPTLSB = SPTLSB
		locationStatus.BotswanaStatus.CGNI = CGNI
		locationStatus.BotswanaStatus.PDF = PDF
		locationStatus.BotswanaStatus.TT140 = TT140
	case "/mnt/ghana":
		locationStatus.GhanaStatus.SE = SE
		locationStatus.GhanaStatus.GL = GL
		locationStatus.GhanaStatus.TXN = TXN
		locationStatus.GhanaStatus.MUL = MUL
		locationStatus.GhanaStatus.VTRAN = VTRAN
		locationStatus.GhanaStatus.VOUT = VOUT
		locationStatus.GhanaStatus.MS = MS
		locationStatus.GhanaStatus.DA = DA
		locationStatus.GhanaStatus.INT00001 = INT00001
		locationStatus.GhanaStatus.INT00003 = INT00003
		locationStatus.GhanaStatus.INT00007 = INT00007
		locationStatus.GhanaStatus.SR = SR
		locationStatus.GhanaStatus.EP = EP
		locationStatus.GhanaStatus.SPTLSB = SPTLSB
		locationStatus.GhanaStatus.CGNI = CGNI
		locationStatus.GhanaStatus.PDF = PDF
		locationStatus.GhanaStatus.TT140 = TT140
	case "/mnt/kenya":
		locationStatus.KenyaStatus.SE = SE
		locationStatus.KenyaStatus.GL = GL
		locationStatus.KenyaStatus.TXN = TXN
		locationStatus.KenyaStatus.MUL = MUL
		locationStatus.KenyaStatus.VTRAN = VTRAN
		locationStatus.KenyaStatus.VOUT = VOUT
		locationStatus.KenyaStatus.MS = MS
		locationStatus.KenyaStatus.DA = DA
		locationStatus.KenyaStatus.INT00001 = INT00001
		locationStatus.KenyaStatus.INT00003 = INT00003
		locationStatus.KenyaStatus.INT00007 = INT00007
		locationStatus.KenyaStatus.SR = SR
		locationStatus.KenyaStatus.EP = EP
		locationStatus.KenyaStatus.SPTLSB = SPTLSB
		locationStatus.KenyaStatus.CGNI = CGNI
		locationStatus.KenyaStatus.PDF = PDF
		locationStatus.KenyaStatus.TT140 = TT140
	case "/mnt/malawi":
		locationStatus.MalawiStatus.SE = SE
		locationStatus.MalawiStatus.GL = GL
		locationStatus.MalawiStatus.TXN = TXN
		locationStatus.MalawiStatus.MUL = MUL
		locationStatus.MalawiStatus.VTRAN = VTRAN
		locationStatus.MalawiStatus.VOUT = VOUT
		locationStatus.MalawiStatus.MS = MS
		locationStatus.MalawiStatus.DA = DA
		locationStatus.MalawiStatus.INT00001 = INT00001
		locationStatus.MalawiStatus.INT00003 = INT00003
		locationStatus.MalawiStatus.INT00007 = INT00007
		locationStatus.MalawiStatus.SR = SR
		locationStatus.MalawiStatus.EP = EP
		locationStatus.MalawiStatus.SPTLSB = SPTLSB
		locationStatus.MalawiStatus.CGNI = CGNI
		locationStatus.MalawiStatus.PDF = PDF
		locationStatus.MalawiStatus.TT140 = TT140
	case "/mnt/namibia":
		locationStatus.NamibiaStatus.SE = SE
		locationStatus.NamibiaStatus.GL = GL
		locationStatus.NamibiaStatus.TXN = TXN
		locationStatus.NamibiaStatus.MUL = MUL
		locationStatus.NamibiaStatus.VTRAN = VTRAN
		locationStatus.NamibiaStatus.VOUT = VOUT
		locationStatus.NamibiaStatus.MS = MS
		locationStatus.NamibiaStatus.DA = DA
		locationStatus.NamibiaStatus.INT00001 = INT00001
		locationStatus.NamibiaStatus.INT00003 = INT00003
		locationStatus.NamibiaStatus.INT00007 = INT00007
		locationStatus.NamibiaStatus.SR = SR
		locationStatus.NamibiaStatus.EP = EP
		locationStatus.NamibiaStatus.SPTLSB = SPTLSB
		locationStatus.NamibiaStatus.CGNI = CGNI
		locationStatus.NamibiaStatus.PDF = PDF
		locationStatus.NamibiaStatus.TT140 = TT140
	case "/mnt/uganda":
		locationStatus.UgandaStatus.SE = SE
		locationStatus.UgandaStatus.GL = GL
		locationStatus.UgandaStatus.TXN = TXN
		locationStatus.UgandaStatus.MUL = MUL
		locationStatus.UgandaStatus.VTRAN = VTRAN
		locationStatus.UgandaStatus.VOUT = VOUT
		locationStatus.UgandaStatus.MS = MS
		locationStatus.UgandaStatus.DA = DA
		locationStatus.UgandaStatus.INT00001 = INT00001
		locationStatus.UgandaStatus.INT00003 = INT00003
		locationStatus.UgandaStatus.INT00007 = INT00007
		locationStatus.UgandaStatus.SR = SR
		locationStatus.UgandaStatus.EP = EP
		locationStatus.UgandaStatus.SPTLSB = SPTLSB
		locationStatus.UgandaStatus.CGNI = CGNI
		locationStatus.UgandaStatus.PDF = PDF
		locationStatus.UgandaStatus.TT140 = TT140
	case "/mnt/ugandadr":
		locationStatus.UgandaDRStatus.SE = SE
		locationStatus.UgandaDRStatus.GL = GL
		locationStatus.UgandaDRStatus.TXN = TXN
		locationStatus.UgandaDRStatus.MUL = MUL
		locationStatus.UgandaDRStatus.VTRAN = VTRAN
		locationStatus.UgandaDRStatus.VOUT = VOUT
		locationStatus.UgandaDRStatus.MS = MS
		locationStatus.UgandaDRStatus.DA = DA
		locationStatus.UgandaDRStatus.INT00001 = INT00001
		locationStatus.UgandaDRStatus.INT00003 = INT00003
		locationStatus.UgandaDRStatus.INT00007 = INT00007
		locationStatus.UgandaDRStatus.SR = SR
		locationStatus.UgandaDRStatus.EP = EP
		locationStatus.UgandaDRStatus.SPTLSB = SPTLSB
		locationStatus.UgandaDRStatus.CGNI = CGNI
		locationStatus.UgandaDRStatus.PDF = PDF
		locationStatus.UgandaDRStatus.TT140 = TT140
	case "/mnt/zambia":
		locationStatus.ZambiaStatus.SE = SE
		locationStatus.ZambiaStatus.GL = GL
		locationStatus.ZambiaStatus.TXN = TXN
		locationStatus.ZambiaStatus.MUL = MUL
		locationStatus.ZambiaStatus.VTRAN = VTRAN
		locationStatus.ZambiaStatus.VOUT = VOUT
		locationStatus.ZambiaStatus.MS = MS
		locationStatus.ZambiaStatus.DA = DA
		locationStatus.ZambiaStatus.INT00001 = INT00001
		locationStatus.ZambiaStatus.INT00003 = INT00003
		locationStatus.ZambiaStatus.INT00007 = INT00007
		locationStatus.ZambiaStatus.SR = SR
		locationStatus.ZambiaStatus.EP = EP
		locationStatus.ZambiaStatus.SPTLSB = SPTLSB
		locationStatus.ZambiaStatus.CGNI = CGNI
		locationStatus.ZambiaStatus.PDF = PDF
		locationStatus.ZambiaStatus.TT140 = TT140
	case "/mnt/zambiadr":
		locationStatus.ZambiaDRStatus.SE = SE
		locationStatus.ZambiaDRStatus.GL = GL
		locationStatus.ZambiaDRStatus.TXN = TXN
		locationStatus.ZambiaDRStatus.MUL = MUL
		locationStatus.ZambiaDRStatus.VTRAN = VTRAN
		locationStatus.ZambiaDRStatus.VOUT = VOUT
		locationStatus.ZambiaDRStatus.MS = MS
		locationStatus.ZambiaDRStatus.DA = DA
		locationStatus.ZambiaDRStatus.INT00001 = INT00001
		locationStatus.ZambiaDRStatus.INT00003 = INT00003
		locationStatus.ZambiaDRStatus.INT00007 = INT00007
		locationStatus.ZambiaDRStatus.SR = SR
		locationStatus.ZambiaDRStatus.EP = EP
		locationStatus.ZambiaDRStatus.SPTLSB = SPTLSB
		locationStatus.ZambiaDRStatus.CGNI = CGNI
		locationStatus.ZambiaDRStatus.PDF = PDF
		locationStatus.ZambiaDRStatus.TT140 = TT140
	case "/mnt/zambiaprod":
		locationStatus.ZambiaProdStatus.SE = SE
		locationStatus.ZambiaProdStatus.GL = GL
		locationStatus.ZambiaProdStatus.TXN = TXN
		locationStatus.ZambiaProdStatus.MUL = MUL
		locationStatus.ZambiaProdStatus.VTRAN = VTRAN
		locationStatus.ZambiaProdStatus.VOUT = VOUT
		locationStatus.ZambiaProdStatus.MS = MS
		locationStatus.ZambiaProdStatus.DA = DA
		locationStatus.ZambiaProdStatus.INT00001 = INT00001
		locationStatus.ZambiaProdStatus.INT00003 = INT00003
		locationStatus.ZambiaProdStatus.INT00007 = INT00007
		locationStatus.ZambiaProdStatus.SR = SR
		locationStatus.ZambiaProdStatus.EP = EP
		locationStatus.ZambiaProdStatus.SPTLSB = SPTLSB
		locationStatus.ZambiaProdStatus.CGNI = CGNI
		locationStatus.ZambiaProdStatus.PDF = PDF
		locationStatus.ZambiaProdStatus.TT140 = TT140

	}
	return nil
}

func (s *service) CreateJSONResponse() LocationStatus {

	return locationStatus
}

func (s *service) ConfirmZimbabweFileAvailability() {
	
	s.ConfirmFileAvailabilityMethod("/mnt/zimbabwe")
}

func (s *service) ConfirmBotswanaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/botswana")
}

func (s *service) ConfirmGhanaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/ghana")
}

func (s *service) ConfirmKenyaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/kenya")
}

func (s *service) ConfirmMalawiFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/malawi")
}

func (s *service) ConfirmNamibiaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/namibia")
}

func (s *service) ConfirmUgandaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/uganda")
}

func (s *service) ConfirmUgandaDRFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/ugandadr")
}

func (s *service) ConfirmZambiaFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/zambia")
}

func (s *service) ConfirmZambiaDRFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/zambiadr")
}

func (s *service) ConfirmZambiaProdFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/zambiaprod")
}

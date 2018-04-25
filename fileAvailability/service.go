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
	GhanaUSDStatus   FilesReceived
	UgandaDRStatus   FilesReceived
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
	EP       bool
	VTRAN    bool
	VOUT     bool
	SPTLSB   bool
	CGNI     bool
	INT00001 bool
	INT00003 bool
	INT00007 bool
	TT140    int
	// PDFs
	VOMTR    bool
	VIFSR    bool
	TIL      bool
	VIMTR    bool
	VOFSR    bool
	MIMTR    bool
	MOMTR    bool
	RR       bool
	MRT      bool
	MOFSR    bool
	MAR      bool
	MIFSR    bool
	DOMTR    bool
	DTIL     bool
}

var locationStatus LocationStatus

func NewService() Service {
	s := &service{}
	//go s.ConfirmBotswanaFileAvailability()
	go s.ConfirmGhanaFileAvailability()
	go s.ConfirmGhanaUSDFileAvailability()
	//go s.ConfirmKenyaFileAvailability()
	go s.ConfirmMalawiFileAvailability()
	go s.ConfirmNamibiaFileAvailability()
	go s.ConfirmUgandaDRFileAvailability()
	go s.ConfirmZambiaProdFileAvailability()
	go s.ConfirmZimbabweFileAvailability()

	return s
}

func (s *service) schedule() {
	confirmBotswanaAvailability := gocron.NewScheduler()
	confirmMalawiAvailability := gocron.NewScheduler()
	confirmGhanaAvailability := gocron.NewScheduler()
	confirmGhanaUSDAvailability := gocron.NewScheduler()
	confirmKenyaAvailability := gocron.NewScheduler()
	confirmNamibiaAvailability := gocron.NewScheduler()
	confirmUgandaDRAvailability := gocron.NewScheduler()
	confirmZambiaProdAvailability := gocron.NewScheduler()
	confirmZimbabweAvailability := gocron.NewScheduler()
	
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
		confirmGhanaAvailability.Every(1).Minute().Do(s.ConfirmGhanaUSDFileAvailability)
		<-confirmGhanaUSDAvailability.Start()
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
		confirmZambiaProdAvailability.Every(1).Minute().Do(s.ConfirmZambiaProdFileAvailability)
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
	TT140 := s.fileCount(path, "TT140")
	// PDFs
	VOMTR, _ := s.pathToMostRecentFile(path, "VISA_OUTGOING_MONET_TRANS_REPORT")
	VIFSR, _ := s.pathToMostRecentFile(path, "VISA_INCOMING_FILES_SUMMARY_REPORT")
	TIL, _ := s.pathToMostRecentFile(path, "TRANS_INPUT_LIST_")
	VIMTR, _ := s.pathToMostRecentFile(path, "VISA_INCOMING_MONET_TRANS_REPORT")
	VOFSR, _ := s.pathToMostRecentFile(path, "VISA_OUTGOING_FILES_SUMMARY_REPORT")
	MIMTR, _ := s.pathToMostRecentFile(path, "MC_INCOMING_MONET_TRANS_REPORT")
	MOMTR, _ := s.pathToMostRecentFile(path, "MC_OUTGOING_MONET_TRANS_REPORT")
	RR, _ := s.pathToMostRecentFile(path, "RECON_REPORT")
	MRT, _ := s.pathToMostRecentFile(path, "MERCH_REJ_TRANS")
	MOFSR, _ := s.pathToMostRecentFile(path, "MC_OUTGOING_FILES_SUMMARY_REPORT")
	MAR, _ := s.pathToMostRecentFile(path, "MASTERCARD_ACKNOWLEDGEMENT_REPORT")
	MIFSR, _ := s.pathToMostRecentFile(path, "MC_INCOMING_FILES_SUMMARY_REPORT")
	DOMTR, _ := s.pathToMostRecentFile(path, "DCI_OUTGOING_MONET_TRANS_REPORT")
	DTIL, _ := s.pathToMostRecentFile(path, "DCI_TRANS_INPUT_LIST")

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
		locationStatus.ZimbabweStatus.TT140 = TT140
		locationStatus.ZimbabweStatus.VOMTR = VOMTR
		locationStatus.ZimbabweStatus.VIFSR = VIFSR
		locationStatus.ZimbabweStatus.TIL = TIL
		locationStatus.ZimbabweStatus.VIMTR = VIMTR
		locationStatus.ZimbabweStatus.VOFSR = VOFSR
		locationStatus.ZimbabweStatus.MIMTR = MIMTR
		locationStatus.ZimbabweStatus.MOMTR = MOMTR
		locationStatus.ZimbabweStatus.RR = RR
		locationStatus.ZimbabweStatus.MRT = MRT
		locationStatus.ZimbabweStatus.MOFSR = MOFSR
		locationStatus.ZimbabweStatus.MAR = MAR
		locationStatus.ZimbabweStatus.MIFSR = MIFSR
		locationStatus.ZimbabweStatus.DOMTR = DOMTR
		locationStatus.ZimbabweStatus.DTIL = DTIL
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
		locationStatus.BotswanaStatus.TT140 = TT140
		locationStatus.BotswanaStatus.VOMTR = VOMTR
		locationStatus.BotswanaStatus.VIFSR = VIFSR
		locationStatus.BotswanaStatus.TIL = TIL
		locationStatus.BotswanaStatus.VIMTR = VIMTR
		locationStatus.BotswanaStatus.VOFSR = VOFSR
		locationStatus.BotswanaStatus.MIMTR = MIMTR
		locationStatus.BotswanaStatus.MOMTR = MOMTR
		locationStatus.BotswanaStatus.RR = RR
		locationStatus.BotswanaStatus.MRT = MRT
		locationStatus.BotswanaStatus.MOFSR = MOFSR
		locationStatus.BotswanaStatus.MAR = MAR
		locationStatus.BotswanaStatus.MIFSR = MIFSR
		locationStatus.BotswanaStatus.DOMTR = DOMTR
		locationStatus.BotswanaStatus.DTIL = DTIL
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
		locationStatus.GhanaStatus.TT140 = TT140
		locationStatus.GhanaStatus.VOMTR = VOMTR
		locationStatus.GhanaStatus.VIFSR = VIFSR
		locationStatus.GhanaStatus.TIL = TIL
		locationStatus.GhanaStatus.VIMTR = VIMTR
		locationStatus.GhanaStatus.VOFSR = VOFSR
		locationStatus.GhanaStatus.MIMTR = MIMTR
		locationStatus.GhanaStatus.MOMTR = MOMTR
		locationStatus.GhanaStatus.RR = RR
		locationStatus.GhanaStatus.MRT = MRT
		locationStatus.GhanaStatus.MOFSR = MOFSR
		locationStatus.GhanaStatus.MAR = MAR
		locationStatus.GhanaStatus.MIFSR = MIFSR
		locationStatus.GhanaStatus.DOMTR = DOMTR
		locationStatus.GhanaStatus.DTIL = DTIL
	case "/mnt/ghanausd":
		locationStatus.GhanaUSDStatus.SE = SE
		locationStatus.GhanaUSDStatus.GL = GL
		locationStatus.GhanaUSDStatus.TXN = TXN
		locationStatus.GhanaUSDStatus.MUL = MUL
		locationStatus.GhanaUSDStatus.VTRAN = VTRAN
		locationStatus.GhanaUSDStatus.VOUT = VOUT
		locationStatus.GhanaUSDStatus.MS = MS
		locationStatus.GhanaUSDStatus.DA = DA
		locationStatus.GhanaUSDStatus.INT00001 = INT00001
		locationStatus.GhanaUSDStatus.INT00003 = INT00003
		locationStatus.GhanaUSDStatus.INT00007 = INT00007
		locationStatus.GhanaUSDStatus.SR = SR
		locationStatus.GhanaUSDStatus.EP = EP
		locationStatus.GhanaUSDStatus.SPTLSB = SPTLSB
		locationStatus.GhanaUSDStatus.CGNI = CGNI
		locationStatus.GhanaUSDStatus.TT140 = TT140
		locationStatus.GhanaUSDStatus.VOMTR = VOMTR
		locationStatus.GhanaUSDStatus.VIFSR = VIFSR
		locationStatus.GhanaUSDStatus.TIL = TIL
		locationStatus.GhanaUSDStatus.VIMTR = VIMTR
		locationStatus.GhanaUSDStatus.VOFSR = VOFSR
		locationStatus.GhanaUSDStatus.MIMTR = MIMTR
		locationStatus.GhanaUSDStatus.MOMTR = MOMTR
		locationStatus.GhanaUSDStatus.RR = RR
		locationStatus.GhanaUSDStatus.MRT = MRT
		locationStatus.GhanaUSDStatus.MOFSR = MOFSR
		locationStatus.GhanaUSDStatus.MAR = MAR
		locationStatus.GhanaUSDStatus.MIFSR = MIFSR
		locationStatus.GhanaUSDStatus.DOMTR = DOMTR
		locationStatus.GhanaUSDStatus.DTIL = DTIL
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
		locationStatus.KenyaStatus.TT140 = TT140
		locationStatus.KenyaStatus.VOMTR = VOMTR
		locationStatus.KenyaStatus.VIFSR = VIFSR
		locationStatus.KenyaStatus.TIL = TIL
		locationStatus.KenyaStatus.VIMTR = VIMTR
		locationStatus.KenyaStatus.VOFSR = VOFSR
		locationStatus.KenyaStatus.MIMTR = MIMTR
		locationStatus.KenyaStatus.MOMTR = MOMTR
		locationStatus.KenyaStatus.RR = RR
		locationStatus.KenyaStatus.MRT = MRT
		locationStatus.KenyaStatus.MOFSR = MOFSR
		locationStatus.KenyaStatus.MAR = MAR
		locationStatus.KenyaStatus.MIFSR = MIFSR
		locationStatus.KenyaStatus.DOMTR = DOMTR
		locationStatus.KenyaStatus.DTIL = DTIL
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
		locationStatus.MalawiStatus.TT140 = TT140
		locationStatus.MalawiStatus.VOMTR = VOMTR
		locationStatus.MalawiStatus.VIFSR = VIFSR
		locationStatus.MalawiStatus.TIL = TIL
		locationStatus.MalawiStatus.VIMTR = VIMTR
		locationStatus.MalawiStatus.VOFSR = VOFSR
		locationStatus.MalawiStatus.MIMTR = MIMTR
		locationStatus.MalawiStatus.MOMTR = MOMTR
		locationStatus.MalawiStatus.RR = RR
		locationStatus.MalawiStatus.MRT = MRT
		locationStatus.MalawiStatus.MOFSR = MOFSR
		locationStatus.MalawiStatus.MAR = MAR
		locationStatus.MalawiStatus.MIFSR = MIFSR
		locationStatus.MalawiStatus.DOMTR = DOMTR
		locationStatus.MalawiStatus.DTIL = DTIL
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
		locationStatus.NamibiaStatus.TT140 = TT140
		locationStatus.NamibiaStatus.VOMTR = VOMTR
		locationStatus.NamibiaStatus.VIFSR = VIFSR
		locationStatus.NamibiaStatus.TIL = TIL
		locationStatus.NamibiaStatus.VIMTR = VIMTR
		locationStatus.NamibiaStatus.VOFSR = VOFSR
		locationStatus.NamibiaStatus.MIMTR = MIMTR
		locationStatus.NamibiaStatus.MOMTR = MOMTR
		locationStatus.NamibiaStatus.RR = RR
		locationStatus.NamibiaStatus.MRT = MRT
		locationStatus.NamibiaStatus.MOFSR = MOFSR
		locationStatus.NamibiaStatus.MAR = MAR
		locationStatus.NamibiaStatus.MIFSR = MIFSR
		locationStatus.NamibiaStatus.DOMTR = DOMTR
		locationStatus.NamibiaStatus.DTIL = DTIL
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
		locationStatus.UgandaDRStatus.TT140 = TT140
		locationStatus.UgandaDRStatus.VOMTR = VOMTR
		locationStatus.UgandaDRStatus.VIFSR = VIFSR
		locationStatus.UgandaDRStatus.TIL = TIL
		locationStatus.UgandaDRStatus.VIMTR = VIMTR
		locationStatus.UgandaDRStatus.VOFSR = VOFSR
		locationStatus.UgandaDRStatus.MIMTR = MIMTR
		locationStatus.UgandaDRStatus.MOMTR = MOMTR
		locationStatus.UgandaDRStatus.RR = RR
		locationStatus.UgandaDRStatus.MRT = MRT
		locationStatus.UgandaDRStatus.MOFSR = MOFSR
		locationStatus.UgandaDRStatus.MAR = MAR
		locationStatus.UgandaDRStatus.MIFSR = MIFSR
		locationStatus.UgandaDRStatus.DOMTR = DOMTR
		locationStatus.UgandaDRStatus.DTIL = DTIL
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
		locationStatus.ZambiaProdStatus.TT140 = TT140
		locationStatus.ZambiaProdStatus.VOMTR = VOMTR
		locationStatus.ZambiaProdStatus.VIFSR = VIFSR
		locationStatus.ZambiaProdStatus.TIL = TIL
		locationStatus.ZambiaProdStatus.VIMTR = VIMTR
		locationStatus.ZambiaProdStatus.VOFSR = VOFSR
		locationStatus.ZambiaProdStatus.MIMTR = MIMTR
		locationStatus.ZambiaProdStatus.MOMTR = MOMTR
		locationStatus.ZambiaProdStatus.RR = RR
		locationStatus.ZambiaProdStatus.MRT = MRT
		locationStatus.ZambiaProdStatus.MOFSR = MOFSR
		locationStatus.ZambiaProdStatus.MAR = MAR
		locationStatus.ZambiaProdStatus.MIFSR = MIFSR
		locationStatus.ZambiaProdStatus.DOMTR = DOMTR
		locationStatus.ZambiaProdStatus.DTIL = DTIL

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

func (s *service) ConfirmGhanaUSDFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/ghanausd")
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

func (s *service) ConfirmUgandaDRFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/ugandadr")
}

func (s *service) ConfirmZambiaProdFileAvailability() {

	s.ConfirmFileAvailabilityMethod("/mnt/zambiaprod")
}

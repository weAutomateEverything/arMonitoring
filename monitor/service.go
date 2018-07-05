package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"log"
)

type Service interface {
	StatusResults() map[string]map[string]string
}

type service struct {
	globalStatus []fileChecker.Service
	store        Store
}

func NewService(store Store, fileStore fileChecker.Store) Service {

	s := &service{store: store}

	log.Println("File arrival confirmation commencing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP747", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006", "SPTLSB"}
	backDatedFiles := []string{"GL", "SE", "TXN", "CGNI", "INT00001", "INT00003", "INT00007", "SR00001", "MUL00002", "MUL00004"}
	afterHoursFiles:= []string{ ".001", ".002", ".003", ".004", ".005", ".006", "SPTLSB"}

	//Zimbabwe
	zimbabwe := fileChecker.NewFileChecker(fileStore, "Zimbabwe", "/mnt/zimbabwe", backDatedFiles, append(common)...)
	s.globalStatus = append(s.globalStatus, zimbabwe)
	//Zambia
	zambia := fileChecker.NewFileChecker(fileStore, "Zambia", "/mnt/zambiaprod", backDatedFiles, append(common)...)
	s.globalStatus = append(s.globalStatus, zambia)
	//Ghana
	ghana := fileChecker.NewFileChecker(fileStore, "Ghana", "/mnt/ghana", backDatedFiles, append(common, "MUL")...)
	s.globalStatus = append(s.globalStatus, ghana)
	//GhanaUSD
	ghanausd := fileChecker.NewFileChecker(fileStore, "GhanaUSD", "/mnt/ghanausd", backDatedFiles, append(common)...)
	s.globalStatus = append(s.globalStatus, ghanausd)
	//Botswana
	botswana := fileChecker.NewFileChecker(fileStore, "Botswana", "/mnt/botswana", backDatedFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
	s.globalStatus = append(s.globalStatus, botswana)
	//Namibia
	namibia := fileChecker.NewFileChecker(fileStore, "Namibia", "/mnt/namibia", backDatedFiles, append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI")...)
	s.globalStatus = append(s.globalStatus, namibia)
	//Malawi
	malawi := fileChecker.NewFileChecker(fileStore, "Malawi", "/mnt/malawi", backDatedFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI")...)
	s.globalStatus = append(s.globalStatus, malawi)
	//Kenya
	kenya := fileChecker.NewFileChecker(fileStore, "Kenya", "/mnt/kenya", backDatedFiles, append(common)...)
	s.globalStatus = append(s.globalStatus, kenya)

	resetsched := gocron.NewScheduler()
	globalStateDailySched := gocron.NewScheduler()

	go func() {
		resetsched.Every(1).Day().At("00:01").Do(s.resetValues)
		<-resetsched.Start()
	}()

	go func() {
		globalStateDailySched.Every(1).Day().At("11:55").Do(s.storeGlobalStateDaily)
		<-globalStateDailySched.Start()
	}()

	return s
}

func (s *service) resetValues() {
	log.Println("Global reset initiated")

	for _, loc := range s.globalStatus {
		loc.Reset()
	}
	log.Println("Global reset completed")
}

func (s *service) StatusResults() map[string]map[string]string {

	response := make(map[string]map[string]string)
	for _, loc := range s.globalStatus {

		response[loc.GetLocationName()] = loc.GetValues()
	}

	return response
}

func (s *service) storeGlobalStateDaily() {
	s.store.addGlobalStateDaily(s.StatusResults())
}

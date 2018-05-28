package monitor

import (
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"log"
)

type Service interface {
	StatusResults() map[string]map[string]string
}

type service struct {
	globalStatus []fileChecker.Service
}

func NewService() Service {

	s := &service{}

	log.Println("File arrival confirmation commencing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP747", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006"}

	//Zimbabwe
	zimbabwe := fileChecker.NewFileChecker("Zimbabwe", "/mnt/zimbabwe", append(common)...)
	s.globalStatus = append(s.globalStatus, zimbabwe)

	//Zambia
	zambia := fileChecker.NewFileChecker("Zambia", "/mnt/zambiaprod", append(common)...)
	s.globalStatus = append(s.globalStatus, zambia)

	//Ghana
	ghana := fileChecker.NewFileChecker("Ghana", "/mnt/ghana", append(common, "MUL")...)
	s.globalStatus = append(s.globalStatus, ghana)
	//GhanaUSD
	ghanausd := fileChecker.NewFileChecker("GhanaUSD", "/mnt/ghanausd", append(common)...)
	s.globalStatus = append(s.globalStatus, ghanausd)
	//Botswana
	botswana := fileChecker.NewFileChecker("Botswana", "/mnt/botswana", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
	s.globalStatus = append(s.globalStatus, botswana)
	//Namibia
	namibia := fileChecker.NewFileChecker("Namibia", "/mnt/namibia", append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "SPTLSB_NA_", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI")...)
	s.globalStatus = append(s.globalStatus, namibia)
	//Malawi
	malawi := fileChecker.NewFileChecker("Malawi", "/mnt/malawi", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
	s.globalStatus = append(s.globalStatus, malawi)
	return s
}

func (s *service) StatusResults() map[string]map[string]string {

	response := make(map[string]map[string]string)
	for _, loc := range s.globalStatus {

		response[loc.GetLocationName()] = loc.GetValues()
	}

	return response
}

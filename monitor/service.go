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
	targets map[string]map[string]string
}

var status = &service{targets: make(map[string]map[string]string)}

func NewService() Service {

	sched := gocron.NewScheduler()
	receiveStatusSched := gocron.NewScheduler()

	go func() {
		receiveStatusSched.Every(1).Day().At("00:00").Do(fileChecker.CreateReceiveStat)
		<-sched.Start()
	}()

	go func() {
		sched.Every(5).Minutes().Do(statusCheck)
		<-sched.Start()
	}()

	fileChecker.CreateReceiveStat()

	statusCheck()

	return status
}

func statusCheck() {

	log.Println("File arrival confirmation commencing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP747", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006"}

	//Zimbabwe
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Zimbabwe", "/mnt/zimbabwe", append(common)...)
		status.targets[locName] = locMap
	}()
	//Zambia
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Zambia", "/mnt/zambiaprod", append(common)...)
		status.targets[locName] = locMap
	}()
	//Ghana
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Ghana", "/mnt/ghana", append(common, "MUL")...)
		status.targets[locName] = locMap
	}()
	//GhanaUSD
	go func() {
		locName, locMap := fileChecker.NewFileChecker("GhanaUSD", "/mnt/ghanausd", append(common)...)
		status.targets[locName] = locMap
	}()
	//Botswana
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Botswana", "/mnt/botswana", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
		status.targets[locName] = locMap
	}()
	//Namibia
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Namibia", "/mnt/namibia", append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "SPTLSB_NA_", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_","CGNI")...)
		status.targets[locName] = locMap
	}()
	//Malawi
	go func() {
		locName, locMap := fileChecker.NewFileChecker("Malawi", "/mnt/malawi", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
		status.targets[locName] = locMap
	}()
}

func (s *service) StatusResults() map[string]map[string]string {

	return status.targets
}

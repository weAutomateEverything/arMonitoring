package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"log"
	"sync"
)

type Service interface {
	StatusResults() []fileChecker.Service
}

type service struct {
	targets []fileChecker.Service
}

var status = &service{}

func NewService() Service {

	sched := gocron.NewScheduler()

	go func() {
		sched.Every(5).Minutes().Do(statusCheck)
		<-sched.Start()
	}()

	s := statusCheck()

	return s
}

func statusCheck() *service {
	s := &service{}

	log.Println("File arrival confirmation commencing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006"}

	var wg sync.WaitGroup

	//Zimbabwe
	wg.Add(7)
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Zimbabwe", "/mnt/zimbabwe", append(common)...))
	}()
	//Zambia
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Zambia", "/mnt/zambiaprod", append(common)...))
	}()
	//Ghana
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Ghana", "/mnt/ghana", append(common, "MUL")...))
	}()
	//GhanaUSD
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("GhanaUSD", "/mnt/ghanausd", append(common)...))
	}()
	//Botswana
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Botswana", "/mnt/botswana", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	}()
	//Namibia
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Namibia", "/mnt/namibia", append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "SPTLSB_NA_", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	}()
	//Malawi
	go func() {
		s.targets = append(s.targets,fileChecker.NewFileChecker("Malawi", "/mnt/malawi", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	}()

	status = s
	return s
}

func (s *service) StatusResults() []fileChecker.Service {

	return status.targets
}

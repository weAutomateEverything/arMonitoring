package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"log"
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
		sched.Every(1).Minute().Do(statusCheck)
		<-sched.Start()
	}()

	s := statusCheck()

	return s
}

func statusCheck() *service {
	s := &service{}
	log.Println("testing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006"}
	//Zimbabwe
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/zimbabwe", append(common)...))
	//Zambia
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/zambiaprod", append(common)...))
	//Ghana
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/ghana", append(common, "MUL")...))
	//GhanaUSD
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/ghanausd", append(common)...))
	//Botswana
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/botswana", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	//Uganda
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/ugandadr", append(common)...))
	//Namibia
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/namibia", append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "SPTLSB_NA_", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	//Malawi
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/malawi", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...))
	//Kenya
	s.targets = append(s.targets, fileChecker.NewFileChecker("/mnt/kenya", append(common)...))
	status = s
	return s
}

func (s *service) StatusResults() []fileChecker.Service {

	return status.targets
}

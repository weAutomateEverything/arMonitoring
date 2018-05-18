package monitor

import (
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"log"
	"sync"
	"time"
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

	log.Println("File arrival confirmation commencing")
	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP", "VTRAN", "VOUT", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", ".001", ".002", ".003", ".004", ".005", ".006"}

	var wg sync.WaitGroup

	locationBuf := make(chan fileChecker.Service)
	//Zimbabwe

	wg.Add(8)
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Zimbabwe", "/mnt/zimbabwe", append(common)...)
		wg.Done()
	}()
	//Zambia
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Zambia", "/mnt/zambiaprod", append(common)...)
		wg.Done()
	}()
	//Ghana
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Ghana", "/mnt/ghana", append(common, "MUL")...)
		time.Sleep(10*time.Second)
		wg.Done()
	}()
	//GhanaUSD
	go func() {
		locationBuf <- fileChecker.NewFileChecker("GhanaUSD", "/mnt/ghanausd", append(common)...)
		wg.Done()
	}()
	//Botswana
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Botswana", "/mnt/botswana", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
		wg.Done()
	}()
	//Namibia
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Namibia", "/mnt/namibia", append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "SPTLSB_NA_", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
		wg.Done()
	}()
	//Malawi
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Malawi", "/mnt/malawi", append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_")...)
		wg.Done()
	}()
	//Kenya
	go func() {
		locationBuf <- fileChecker.NewFileChecker("Kenya", "/mnt/kenya", append(common)...)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(locationBuf)
	}()

	for ch := range locationBuf {
		s.targets = append(s.targets, ch)
	}

	status = s
	return s
}

func (s *service) StatusResults() []fileChecker.Service {

	return status.targets
}

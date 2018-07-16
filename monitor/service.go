package monitor

import (

	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

)

type Service interface {
	StatusResults() map[string]map[string]string
	resetValues()
	resetAfterHoursValues()
	storeGlobalStateDaily()
}

type service struct {
	globalStatus []fileChecker.Service
	store        Store
}

func NewService(fieldKeys []string, logger log.Logger, store Store, fileStore fileChecker.Store) Service {

	s := &service{store: store}

	common := []string{"SE", "GL", "TXN", "DA", "MS", "EP747", "VISA_OUTGOING_MONET_TRANS_REPORT", "VISA_INCOMING_FILES_SUMMARY_REPORT", "TRANS_INPUT_LIST_", "VISA_INCOMING_MONET_TRANS_REPORT", "VISA_OUTGOING_FILES_SUMMARY_REPORT", "MC_INCOMING_MONET_TRANS_REPORT", "MC_OUTGOING_MONET_TRANS_REPORT", "RECON_REPORT", "MERCH_REJ_TRANS", "MC_OUTGOING_FILES_SUMMARY_REPORT", "MASTERCARD_ACKNOWLEDGEMENT_REPORT", "MC_INCOMING_FILES_SUMMARY_REPORT", "SPTLSB"}
	backDatedFiles := []string{"GL", "SE", "TXN", "CGNI", "INT00001", "INT00003", "INT00007", "SR00001", "MUL00002", "MUL00004"}
	afterHoursFiles:= []string{ ".001", ".002", ".003", ".004", ".005", ".006", "SPTLSB"}

	//Zimbabwe
	zimbabwe := fileChecker.NewFileChecker(fileStore, "Zimbabwe", "/mnt/zimbabwe", backDatedFiles, afterHoursFiles, append(common, "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	zimbabwe = fileChecker.NewLoggingService(log.With(logger, "component", "zimbabweFileChecker"), zimbabwe)
	zimbabwe = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "zimbabweFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "zimbabweFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), zimbabwe)
	s.globalStatus = append(s.globalStatus, zimbabwe)

	//Zambia
	zambia := fileChecker.NewFileChecker(fileStore, "Zambia", "/mnt/zambiaprod", backDatedFiles, afterHoursFiles, append(common, "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	zambia = fileChecker.NewLoggingService(log.With(logger, "component", "zambiaFileChecker"), zambia)
	zambia = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "zambiaFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "zambiaFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), zambia)
	s.globalStatus = append(s.globalStatus, zambia)

	//Ghana
	ghana := fileChecker.NewFileChecker(fileStore, "Ghana", "/mnt/ghana", backDatedFiles, afterHoursFiles, append(common, "MUL", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	ghana = fileChecker.NewLoggingService(log.With(logger, "component", "ghanaFileChecker"), ghana)
	ghana = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "ghanaFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "ghanaFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), ghana)
	s.globalStatus = append(s.globalStatus, ghana)

	//GhanaUSD
	ghanausd := fileChecker.NewFileChecker(fileStore, "GhanaUSD", "/mnt/ghanausd", backDatedFiles, afterHoursFiles, append(common, "VTRAN", "VOUT")...)
	ghanausd = fileChecker.NewLoggingService(log.With(logger, "component", "ghanausdFileChecker"), ghanausd)
	ghanausd = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "ghanausdFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "ghanausdFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), ghanausd)
	s.globalStatus = append(s.globalStatus, ghanausd)

	//Botswana
	botswana := fileChecker.NewFileChecker(fileStore, "Botswana", "/mnt/botswana", backDatedFiles, afterHoursFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	botswana = fileChecker.NewLoggingService(log.With(logger, "component", "botswanaFileChecker"), botswana)
	botswana = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "botswanaFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "botswanaFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), botswana)
	s.globalStatus = append(s.globalStatus, botswana)

	//Namibia
	namibia := fileChecker.NewFileChecker(fileStore, "Namibia", "/mnt/namibia", backDatedFiles, afterHoursFiles, append(common, "MUL", "INT00001", "INT00003", "INT00007", "SR00001", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	namibia = fileChecker.NewLoggingService(log.With(logger, "component", "namibiaFileChecker"), namibia)
	namibia = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "namibiaFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "namibiaFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), namibia)
	s.globalStatus = append(s.globalStatus, namibia)

	//Malawi
	malawi := fileChecker.NewFileChecker(fileStore, "Malawi", "/mnt/malawi", backDatedFiles, afterHoursFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	malawi = fileChecker.NewLoggingService(log.With(logger, "component", "malawiFileChecker"), malawi)
	malawi = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "malawiFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "malawiFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), malawi)
	s.globalStatus = append(s.globalStatus, malawi)

	//Kenya
	kenya := fileChecker.NewFileChecker(fileStore, "Kenya", "/mnt/kenya", backDatedFiles, afterHoursFiles, append(common)...)
	kenya = fileChecker.NewLoggingService(log.With(logger, "component", "kenyaFileChecker"), kenya)
	kenya = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "kenyaFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "kenyaFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), kenya)
	s.globalStatus = append(s.globalStatus, kenya)

	//Lesotho
	lesotho := fileChecker.NewFileChecker(fileStore, "Lesotho", "/mnt/lesotho", backDatedFiles, afterHoursFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	lesotho = fileChecker.NewLoggingService(log.With(logger, "component", "lesothoFileChecker"), lesotho)
	lesotho = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "lesothoFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "lesothoFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), lesotho)
	s.globalStatus = append(s.globalStatus, lesotho)

	//Swaziland
	swaziland := fileChecker.NewFileChecker(fileStore, "Swaziland", "/mnt/swaziland", backDatedFiles, afterHoursFiles, append(common, "MUL", "DCI_OUTGOING_MONET_TRANS_REPORT", "DCI_TRANS_INPUT_LIST_", "CGNI", "VTRAN", "VOUT", ".001", ".002", ".003", ".004", ".005", ".006")...)
	swaziland = fileChecker.NewLoggingService(log.With(logger, "component", "swazilandFileChecker"), swaziland)
	swaziland = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "swazilandFileChecker",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "swazilandFileChecker",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), swaziland)
	s.globalStatus = append(s.globalStatus, swaziland)

	resetsched := gocron.NewScheduler()
	afterHoursResetsched := gocron.NewScheduler()
	globalStateDailySched := gocron.NewScheduler()

	go func() {
		resetsched.Every(1).Day().At("00:01").Do(s.resetValues)
		<-resetsched.Start()
	}()

	go func() {
		afterHoursResetsched.Every(1).Day().At("17:00").Do(s.resetAfterHoursValues)
		<-afterHoursResetsched.Start()
	}()

	go func() {
		globalStateDailySched.Every(1).Day().At("23:55").Do(s.storeGlobalStateDaily)
		<-globalStateDailySched.Start()
	}()

	return s
}

func (s *service) resetValues() {
	for _, loc := range s.globalStatus {
		loc.Reset()
	}
}

func (s *service) resetAfterHoursValues() {
	for _, loc := range s.globalStatus {
		loc.ResetAfterHours()
	}
}

func (s *service) StatusResults() map[string]map[string]string {

	response := make(map[string]map[string]string)
	for _, loc := range s.globalStatus {

		response[loc.GetLocationName()] = loc.GetValues()
	}

	return response
}

func (s *service) storeGlobalStateDaily() {
	s.store.setGlobalStateDaily(s.StatusResults())
}


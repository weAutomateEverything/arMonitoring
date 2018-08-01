package monitor

import (
	"github.com/go-kit/kit/log"
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"github.com/weAutomateEverything/fileMonitorService/jsonFileInteraction"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"strings"
)

type Service interface {
	StatusResults() map[string]map[string]string
	getDatedGlobalStateDaily(date string) (map[string]map[string]string, error)
	resetValues()
	resetAfterHoursValues()
	storeGlobalStateDaily()
}

type service struct {
	globalStatus []fileChecker.Service
	store        Store
}

func NewService(json jsonFileInteraction.Service, fieldKeys []string, logger log.Logger, store Store, fileStore fileChecker.Store) Service {

	s := &service{store: store}

	locations := json.ReturnLocationsArray()

	for index := range locations {

		location := fileChecker.NewFileChecker(json, fileStore, strings.Title(locations[index].Name), locations[index].MountPath, json.ReturnBackdatedFilesArray(), json.ReturnAfterHoursFilesArray(), append(json.ReturnCommonFilesArray(), locations[index].Files...)...)
		location = fileChecker.NewLoggingService(log.With(logger, "component", locations[index].Name+"FileChecker"), location)
		location = fileChecker.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: locations[index].Name + "FileChecker",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: locations[index].Name + "FileChecker",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys), location)
		s.globalStatus = append(s.globalStatus, location)
	}

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

func (s *service) getDatedGlobalStateDaily(date string) (map[string]map[string]string, error) {
	return s.store.getGlobalStateDailyForThisDate(date)
}

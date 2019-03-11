package monitor

import (
	"github.com/go-kit/kit/log"
	"github.com/jasonlvhit/gocron"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"github.com/weAutomateEverything/fileMonitorService/jsonFileInteraction"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"strings"
	"github.com/weAutomateEverything/fileMonitorService/cyberArk"
	"time"
)

type Service interface {
	StatusResults() Response
	getDatedGlobalStateDaily(date string) (Response, error)
	resetValues()
	resetAfterHoursValues()
	storeGlobalStateDaily()
	//updateCyberarkCredentials()
}

type service struct {
	globalStatus []fileChecker.Service
	store        Store
	cark 		 cyberArk.Service
}

type Response struct {
	Locations []Location `json:"locations"`
}

type Location struct {
	Tab          string `json:"tab"`
	LocationName string `json:"locationname"`
	Date 		 string `json:"date"`
	Files        map[string]string `json:"files"`
}

//Create new Filechecker instance in memory for each location
func NewService(json jsonFileInteraction.Service, fieldKeys []string, logger log.Logger, store Store, fileStore fileChecker.Store) Service {

	s := &service{store: store}

	locations := json.ReturnLocationsArray()

	for index := range locations {

		location := fileChecker.NewFileChecker(json,locations[index].TabNumber ,fileStore, strings.Title(locations[index].Name), locations[index].MountPath, json.ReturnBackdatedFilesArray(), json.ReturnAfterHoursFilesArray(), locations[index].Files...)
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

func (s *service) StatusResults() Response {

	response := Response{}
	for _, loc := range s.globalStatus {
		response.Locations = append(response.Locations, Location{loc.GetTabNumber(),loc.GetLocationName(),time.Now().Format("20060102"),loc.GetValues()})
	}

	return response
}

func (s *service) storeGlobalStateDaily() {
	s.store.setGlobalStateDaily(s.StatusResults().Locations)
}

func (s *service) getDatedGlobalStateDaily(date string) (Response, error) {
	return s.store.getGlobalStateDailyForThisDate(date)
}

//func (s *service) updateCyberarkCredentials() {
//	s.cark.GetCyberarkPassword()
//}

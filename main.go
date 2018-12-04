package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"

	"github.com/weAutomateEverything/fileMonitorService/cyberArk"
	"github.com/weAutomateEverything/fileMonitorService/fileChecker"
	"github.com/weAutomateEverything/fileMonitorService/jsonFileInteraction"
	"github.com/weAutomateEverything/fileMonitorService/monitor"
	"github.com/weAutomateEverything/fileMonitorService/snmp"
	"github.com/weAutomateEverything/fileMonitorService/database"
	"os/signal"
	"syscall"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {

	var logger log.Logger
	var db *mgo.Database
	var errorage error
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	go func() {
		sn := snmp.NewService()
		sn = snmp.NewLoggingService(log.With(logger, "component", "snmp"), sn)
	}()

	fieldKeys := []string{"method"}


	for true {
		db, errorage = database.NewConnection()
		if errorage != nil {
			fmt.Println("Failed to access MongoDB. Retrying in 60 seconds")
			time.Sleep(1 * time.Minute)
		}else {
			break
		}
	}

	cark := cyberArk.NewCyberarkRetreivalService()

	dailyStore := monitor.NewMongoStore(db)
	recentStore := fileChecker.NewMongoStore(db)

	json := jsonFileInteraction.NewJSONService()
	json = jsonFileInteraction.NewLoggingService(log.With(logger, "component", "jsonFileInteraction"), json)
	json = jsonFileInteraction.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "jsonFileInteraction",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "jsonFileInteraction",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), json)

	mon := monitor.NewService(cark, json, fieldKeys, logger, dailyStore, recentStore)
	mon = monitor.NewLoggingService(log.With(logger, "component", "fileMonitor"), mon)
	mon = monitor.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "MonitorService",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "MonitorService",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), mon)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/fileStatus", monitor.MakeHandler(mon, httpLogger, nil))
	mux.Handle("/setGlobalState", monitor.MakeHandler(mon, httpLogger, nil))
	mux.Handle("/updateCredentials", monitor.MakeHandler(mon, httpLogger, nil))
	mux.Handle("/backdated", monitor.MakeHandler(mon, httpLogger, nil))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":8002", "msg", "listening")
		errs <- http.ListenAndServe(":8002", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

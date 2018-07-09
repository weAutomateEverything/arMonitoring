package fileChecker

import (
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) GetValues() map[string]string {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetValues").Add(1)
		s.requestLatency.With("method", "GetValues").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetValues()
}

func (s *instrumentingService) GetLocationName() string {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetLocationName").Add(1)
		s.requestLatency.With("method", "GetLocationName").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetLocationName()
}

func (s *instrumentingService) Reset() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Reset").Add(1)
		s.requestLatency.With("method", "Reset").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.Reset()
}

func (s *instrumentingService) ResetAfterHours(){
	defer func(begin time.Time) {
		s.requestCount.With("method", "ResetAfterHours").Add(1)
		s.requestLatency.With("method", "ResetAfterHours").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ResetAfterHours()
}

func (s *instrumentingService) setValues(name, mountpath string, bdFiles []string, files []string, store Store) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Reset").Add(1)
		s.requestLatency.With("method", "Reset").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.setValues(name, mountpath, bdFiles, files, store)
}

func (s *instrumentingService) storeLocationStateRecent(name string, fileStatus map[string]string){
	defer func(begin time.Time) {
		s.requestCount.With("method", "ResetAfterHours").Add(1)
		s.requestLatency.With("method", "ResetAfterHours").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.storeLocationStateRecent(name, fileStatus)
}

func (s *instrumentingService) setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Reset").Add(1)
		s.requestLatency.With("method", "Reset").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.setFileStatus(name, dirPath, fileContains, bdFiles,store)
}

func (s *instrumentingService) 	getListOfFilesInPath(path string) ([]string, error){
	defer func(begin time.Time) {
		s.requestCount.With("method", "ResetAfterHours").Add(1)
		s.requestLatency.With("method", "ResetAfterHours").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.getListOfFilesInPath(path)
}
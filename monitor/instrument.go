package monitor

import (
	"github.com/go-kit/kit/metrics"
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

//func (s *instrumentingService) ConfirmZimbabweFileAvailability() {
//	defer func(begin time.Time) {
//		s.requestCount.With("method", "ConfirmUgandaFileAvailability").Add(1)
//		s.requestLatency.With("method", "ConfirmUgandaFileAvailability").Observe(time.Since(begin).Seconds())
//	}(time.Now())
//	s.Service.ConfirmZimbabweFileAvailability()
//}
//
//func (s *instrumentingService) GetFilesInPath(path string) ([]string, error) {
//	defer func(begin time.Time) {
//		s.requestCount.With("method", "GetFilesInPath").Add(1)
//		s.requestLatency.With("method", "GetFilesInPath").Observe(time.Since(begin).Seconds())
//	}(time.Now())
//	return s.Service.GetFilesInPath(path)
//}

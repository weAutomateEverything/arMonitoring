package fileAvailability

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

func (s *instrumentingService) ConfirmUgandaFileAvailability() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetFilesInPath").Add(1)
		s.requestLatency.With("method", "GetFilesInPath").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ConfirmUgandaFileAvailability()
}

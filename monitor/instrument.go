package monitor

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

func (s *instrumentingService) StatusResults() Response{
	defer func(begin time.Time) {
		s.requestCount.With("method", "StatusResults").Add(1)
		s.requestLatency.With("method", "StatusResults").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.StatusResults()
}

func (s *instrumentingService) resetValues() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "resetValues").Add(1)
		s.requestLatency.With("method", "resetValues").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.resetValues()
}

func (s *instrumentingService) resetAfterHoursValues() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "resetAfterHoursValues").Add(1)
		s.requestLatency.With("method", "resetAfterHoursValues").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.resetAfterHoursValues()
}

func (s *instrumentingService) storeGlobalStateDaily() {
	defer func(begin time.Time) {
		s.requestCount.With("method", "storeGlobalStateDaily").Add(1)
		s.requestLatency.With("method", "storeGlobalStateDaily").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.storeGlobalStateDaily()
}


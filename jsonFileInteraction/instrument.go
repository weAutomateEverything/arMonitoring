package jsonFileInteraction

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

func (s *instrumentingService) ReturnFileNamesArray() []FileName {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ReturnFileNamesArray").Add(1)
		s.requestLatency.With("method", "ReturnFileNamesArray").Observe(time.Since(begin).Seconds())
	}(time.Now())
	array := s.Service.ReturnFileNamesArray()
	return array
}

func (s *instrumentingService) UnmarshalJSONFile(file string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "UnmarshalJSONFile").Add(1)
		s.requestLatency.With("method", "UnmarshalJSONFile").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.UnmarshalJSONFile(file)
}

package monitor

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) StatusResults() Response{
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "StatusResults",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.StatusResults()
}

func (s *loggingService) resetValues() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "resetValues",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.resetValues()
}

func (s *loggingService) resetAfterHoursValues() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "resetAfterHoursValues",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.resetAfterHoursValues()
}

func (s *loggingService) storeGlobalStateDaily() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "storeGlobalStateDaily",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.storeGlobalStateDaily()
}
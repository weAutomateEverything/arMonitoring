package fileAvailability

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

func (s *loggingService) ConfirmUgandaFileAvailability() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ConfirmUgandaFileAvailability",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ConfirmZimbabweFileAvailability()
}

func (s *loggingService) GetFilesInPath(path string) ([]string, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ConfirmUgandaFileAvailability",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetFilesInPath(path)
}

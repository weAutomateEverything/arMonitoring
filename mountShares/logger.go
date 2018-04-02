package mountShares

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

func (s *loggingService) MountShares() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "mountShares",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.MountShares()
}

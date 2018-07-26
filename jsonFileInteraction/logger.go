package jsonFileInteraction

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

func (s *loggingService) ReturnFileNamesArray() []FileName {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ReturnFileNamesArray",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.ReturnFileNamesArray()
}

func (s *loggingService) UnmarshalJSONFile(file string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "UnmarshalJSONFile",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.UnmarshalJSONFile(file)
}

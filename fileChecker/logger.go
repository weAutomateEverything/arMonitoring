package fileChecker

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

func (s *loggingService) GetValues() map[string]string{
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetValues",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetValues()
}

func (s *loggingService) GetLocationName() string {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetLocationName",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetLocationName()
}

func (s *loggingService) Reset() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Reset",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.Reset()
}

func (s *loggingService) ResetAfterHours() {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ResetAfterHours",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.ResetAfterHours()
}

func (s *loggingService) setValues(name, mountpath string, bdFiles []string, files []string, store Store){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "setValues",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.setValues(name, mountpath, bdFiles, files, store)
}

func (s *loggingService) storeLocationStateRecent(name string, fileStatus map[string]string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "storeLocationStateRecent",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.storeLocationStateRecent(name, fileStatus)
}

func (s *loggingService) setFileStatus(name, dirPath, fileContains string, bdFiles []string, store Store) (string, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "setFileStatus",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.setFileStatus(name, dirPath, fileContains, bdFiles, store)
}

func (s *loggingService) getListOfFilesInPath(path string) ([]string, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "getListOfFilesInPath",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.getListOfFilesInPath(path)
}

func (s *loggingService) convertFileNamesToHumanReadableNames() map[string]string {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "convertFileNamesToHumanReadableNames",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.convertFileNamesToHumanReadableNames()
}
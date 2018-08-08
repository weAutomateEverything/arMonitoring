package snmp

import (
	"github.com/go-kit/kit/log"
	g "github.com/soniah/gosnmp"
	"time"
	"net"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) startSnmpServer(){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "startSnmpServer",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.startSnmpServer()
}

func (s *loggingService) handleTrap(packet *g.SnmpPacket, addr *net.UDPAddr) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "handleTrap",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.handleTrap(packet,addr)
}
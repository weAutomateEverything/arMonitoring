package snmp

import (
"net"
"bytes"
"fmt"
g "github.com/soniah/gosnmp"
"log"
"os"
)

type Service interface {
	startSnmpServer()
	handleTrap(packet *g.SnmpPacket, addr *net.UDPAddr)
}

type service struct {
}

func NewService() Service{
	s := &service{}
	go func() {
		s.startSnmpServer()
	}()

	return s
}

func (s *service) startSnmpServer() {
	log.Println("Starting SNMP server. Listening on port 8003")
	tl := g.NewTrapListener()
	tl.OnNewTrap = s.handleTrap
	tl.Params = g.Default
	tl.Params.Logger = log.New(os.Stdout, "", 0)
	err := tl.Listen("0.0.0.0:8003")
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}

}

func (s service) handleTrap(packet *g.SnmpPacket, addr *net.UDPAddr) {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("got trapdata from %s\n", addr.IP))
	b.WriteString("\n")
	for _, v := range packet.Variables {
		switch v.Type {
		case g.OctetString:
			c := v.Value.([]byte)
			b.WriteString(fmt.Sprintf("OID: %s, string: %x\n", v.Name, c))
			b.WriteString("\n")
		default:
			b.WriteString(fmt.Sprintf("trap: %+v\n", v))
			b.WriteString("\n")

		}
	}

	fmt.Println(b)
}
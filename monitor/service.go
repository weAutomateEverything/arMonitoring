package monitor

import "github.com/weAutomateEverything/fileMonitorService/fileChecker"

type Service interface {

}

type service struct {
	targets []fileChecker.Service
}

func NewService() {
	s := &service{}

	common := []string{"SE","GL","TXN","MUL","VTRAN","VOUT"}
	//ZImbabwe
	s.targets = append(s.targets,fileChecker.NewFileChecker("/mnt/zimbabwe",append(common,"xxxx")...))
	//Zambia
	s.targets = append(s.targets,fileChecker.NewFileChecker("/mnt/zambia",append(common,"xxxx")...))



}

package mountShares

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"github.com/weAutomateEverything/arMonitoring/fileAvailability"
)

type Service interface {
	MountShares()
	ConfirmSharesExist()
}

type service struct {
	fileAvailable fileAvailability.Service
	SharesAvailable bool
}

func NewService(files fileAvailability.Service) Service {
	s := &service{fileAvailable: files}

	return s
}

func (s *service) ConfirmSharesExist(){
	s.fileAvailable.
}

func (s *service) MountShares() {

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go mountCommand("/mnt/share01", getShare01_Location())
	go mountCommand("/mnt/share02", getShare02_Location())
}

func mountCommand(mount, share string) []byte {
	cmd, err := exec.Command("sh", "-c", fmt.Sprintf("'%v %v cifs uid=0,gid=0,user=%v,password=%v 0 0' | tee -a /etc/fstab /dev/null", share, mount, getShare_User(), getShare_Password())).Output()
	if err != nil {
		log.Println(err)
	}
	return cmd
}

func getShare01_Location() string {
	return os.Getenv("SHARE01")
}

func getShare02_Location() string {
	return os.Getenv("SHARE02")
}

func getShare_User() string {
	return os.Getenv("SHARE_USER")
}
func getShare_Password() string {
	return os.Getenv("SHARE_PASSWORD")
}

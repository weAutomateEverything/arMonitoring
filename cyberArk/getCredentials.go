package cyberArk

import (
	"github.com/hoop33/go-cyberark"
	"io/ioutil"
	"log"
	"os"
)

type Service interface {
	GetCyberarkPassword() error
}

type service struct {
	Username string
	Password string
}

func NewCyberarkRetreivalService() Service {
	s := &service{}
	if err := s.GetCyberarkPassword(); err != nil {
		log.Println("Failed to access Cyberark vault")
	}
	return s
}

func (s *service) GetCyberarkPassword() error {
	client, err := cyberark.NewClient(
		cyberark.SetHost("https://epvs.za.sbicdirectory.com"),
	)
	if err != nil {
		log.Println(err.Error())
	}

	ret, err := client.GetPassword().
		AppID(getCyberarkAppID()).
		Safe(getCyberarkSafe()).
		Object(getCyberarkObject()).
		UserName(getCyberarkUsername()).
		Do()
	if err != nil {
		return err
	}

	if ret.ErrorCode != "" {
		log.Println(ret.ErrorCode)
	}

	s.Username = ret.UserName
	s.Password = ret.Content

	log.Println("Successfully retreived credentials from Cyberark vault")

	s.updateCedentialsFile()

	return nil
}

func (s *service) updateCedentialsFile() {

	err := ioutil.WriteFile("/opt/app/creds", []byte("bolloks"), 0755)
	if err != nil {
		log.Printf("Failed writing to credentials file with the following error: %v", err)
	}

}

func getCyberarkAppID() string {
	env := os.Getenv("CYBERARK_APPID")
	return env
}

func getCyberarkSafe() string {
	env := os.Getenv("CYBERARK_SAFE")
	return env
}

func getCyberarkObject() string {
	env := os.Getenv("CYBERARK_OBJECT")
	return env
}

func getCyberarkUsername() string {
	env := os.Getenv("CYBERARK_USERNAME")
	return env
}

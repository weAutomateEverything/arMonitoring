package cyberArk

import (
	"github.com/hoop33/go-cyberark"
	"log"
	"os"
)

type Service interface {
}

type service struct {
	Username string
	Password string
}

func NewCyberarkRetreivalService() Service {
	s := service{}
	s.GetCyberarkPassword()
	return s
}

func (s *service) GetCyberarkPassword() {
	client, err := cyberark.NewClient(
		cyberark.SetHost("https://epvs.za.sbicdirectory.com"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	ret, err := client.GetPassword().
		AppID(getCyberarkAppID()).
		Safe(getCyberarkSafe()).
		Object(getCyberarkObject()).
		UserName(getCyberarkUsername()).
		Do()
	if err != nil {
		log.Println(err.Error())
	}

	if ret.ErrorCode != "" {
		log.Println(ret.ErrorCode)
	}

	s.Username = ret.Content
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

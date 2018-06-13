package monitor

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type Store interface {
	addGlobalStateDaily(globalFileStatus map[string]map[string]string)
}

type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	LastUpdate       time.Time
	GlobalFileStatus map[string]map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) addGlobalStateDaily(globalFileStatus map[string]map[string]string) {
	log.Println("Storing daily global state")
	c := s.mongo.C("GlobalStateDaily")
	stateItem := globalState{LastUpdate: time.Now(), GlobalFileStatus: globalFileStatus}
	err := c.Insert(stateItem)
	if err != nil {
		log.Println("Failed to insert global state data with the following error: ", err)
	}
}

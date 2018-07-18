package monitor

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type Store interface {
	setGlobalStateDaily(globalFileStatus map[string]map[string]string) error
	getGlobalStateDailyForThisDate(searchDate string) (map[string]map[string]string, error)
}

type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	Date             string `bson:"_id,omitempty"`
	GlobalFileStatus map[string]map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) setGlobalStateDaily(globalFileStatus map[string]map[string]string) error {
	log.Println("Storing daily global state")
	c := s.mongo.C("GlobalStateDaily")
	stateItem := globalState{Date: time.Now().Format("02012006"), GlobalFileStatus: globalFileStatus}
	return c.Insert(stateItem)
}

func (s mongoStore) getGlobalStateDailyForThisDate(searchDate string) (map[string]map[string]string, error) {
	log.Printf("Retreiving global state for %v", searchDate)
	c := s.mongo.C("GlobalStateDaily")
	var gs globalState
	err := c.FindId(searchDate).One(&gs)
	if err != nil {
		return nil, err
	}
	return gs.GlobalFileStatus, nil
}

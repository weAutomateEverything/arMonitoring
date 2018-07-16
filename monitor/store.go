package monitor

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	setGlobalStateDaily(globalFileStatus map[string]map[string]string) error
	getGlobalStateDailyForThisDate(searchDate string) *globalState
}

type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	Date       time.Time `bson:"_id,omitempty"`
	GlobalFileStatus map[string]map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) setGlobalStateDaily(globalFileStatus map[string]map[string]string) error {
	log.Println("Storing daily global state")
	c := s.mongo.C("GlobalStateDaily")
	stateItem := globalState{Date: time.Now(), GlobalFileStatus: globalFileStatus}
	return c.Insert(stateItem)
}

func (s mongoStore) getGlobalStateDailyForThisDate(searchDate string) *globalState {
	log.Printf("Retreiving global state for %v", searchDate)
	c := s.mongo.C("GlobalStateDaily")
	globalStateItem := &globalState{}
	err := c.Find(bson.M{"Date": searchDate}).One(&globalStateItem)
	if err != nil {
		log.Printf("Faild to retreive Global State for %v", searchDate)
	}
	return globalStateItem
}
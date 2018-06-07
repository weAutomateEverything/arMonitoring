package monitor

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Store interface {
	getGlobalStateDaily() map[string]map[string]string
	addGlobalStateDaily(globalFileStatus map[string]map[string]string) error
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

func (s mongoStore) getGlobalStateDaily() map[string]map[string]string {
	c := s.mongo.C("GlobalStateDaily")
	var gloState globalState
	c.Find(bson.M{}).One(&gloState)
	return gloState.GlobalFileStatus
}

func (s mongoStore) addGlobalStateDaily(globalFileStatus map[string]map[string]string) error {
	log.Println("Storing daily global state")
	c := s.mongo.C("GlobalStateDaily")
	stateItem := globalState{}
	stateItem.LastUpdate = time.Now()
	stateItem.GlobalFileStatus = globalFileStatus
	err := c.Insert(stateItem)
	if err != nil {
		return err
	}
	return nil
}

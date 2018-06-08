package fileChecker

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Store interface {
	getLocationStateRecent(locationName string) map[string]string
	addLocationStateRecent(locationName string, locationFileStatus map[string]string) error
}

type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	LastUpdate         time.Time
	LocationName       string
	LocationFileStatus map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) getLocationStateRecent(locationName string) map[string]string {
	c := s.mongo.C("GlobalStateRecent")
	var gloState globalState
	c.Find(bson.M{"locationname": locationName}).One(&gloState)
	return gloState.LocationFileStatus
}

func (s mongoStore) addLocationStateRecent(locationName string, locationFileStatus map[string]string) error {
	log.Println("Storing recent global state")
	c := s.mongo.C("GlobalStateRecent")
	stateItem := globalState{}
	stateItem.LastUpdate = time.Now()
	stateItem.LocationName = locationName
	stateItem.LocationFileStatus = locationFileStatus
	err := c.Insert(stateItem)
	if err != nil {
		return err
	}
	return nil
}

package fileChecker

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Store interface {
	getLocationStateRecent(locationName string) (map[string]string, error)
	setLocationStateRecent(locationName string, locationFileStatus map[string]string) error
}

type mongoStore struct {
	mongo *mgo.Database
}

type locationState struct {
	LocationName       string `bson:"_id,omitempty"`
	LocationFileStatus map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) getLocationStateRecent(locationName string) (map[string]string, error) {
	c := s.mongo.C("LocationStateRecent")
	var locState locationState
	err := c.FindId(locationName).One(&locState)
	if err != nil {
		return nil, err
	}
	return locState.LocationFileStatus, nil
}


func (s mongoStore) setLocationStateRecent(locationName string, locationFileStatus map[string]string) error {
	log.Println("Storing/updating recent location state")
	c := s.mongo.C("LocationStateRecent")

	var locState locationState
	err := c.FindId(locationName).One(&locState)
	if err == nil {
		locState.LocationFileStatus = locationFileStatus
		return c.UpdateId(locationName, locState.LocationFileStatus)
	}
	log.Print(err)
	state := locationState{LocationName: locationName, LocationFileStatus: locationFileStatus}
	return c.Insert(state)
}


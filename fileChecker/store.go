package fileChecker

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Store interface {
	getLocationStateRecent(locationName string) map[string]string
	setLocationStateRecent(locationName string, locationFileStatus map[string]string)
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

func (s mongoStore) getLocationStateRecent(locationName string) map[string]string {
	c := s.mongo.C("LocationStateRecent")
	var locState locationState
	c.FindId(locationName).One(&locState)
	return locState.LocationFileStatus
}

func (s mongoStore) setLocationStateRecent(locationName string, locationFileStatus map[string]string) {
	log.Println("Storing/updating recent location state")
	c := s.mongo.C("LocationStateRecent")

	var locState locationState
	err := c.FindId(locationName).One(&locState)
	if err == nil {
		locState.LocationFileStatus = locationFileStatus
		c.UpdateId(locationName, locState.LocationFileStatus)
	} else {
		log.Print(err)
		state := locationState{LocationName: locationName, LocationFileStatus: locationFileStatus}
		err := c.Insert(state)
		if err != nil {
			log.Println("Failed to insert location state data with the following error: ", err)
		}
	}
}

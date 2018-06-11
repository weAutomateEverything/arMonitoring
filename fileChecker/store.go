package fileChecker

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Store interface {
	getLocationStateRecent(locationName string) map[string]string
	setLocationStateRecent(locationName string, locationFileStatus map[string]string)
}

type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	LocationName        string `bson:"_id,omitempty"`
	LocationFileStatus map[string]string
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) getLocationStateRecent(locationName string) map[string]string {
	c := s.mongo.C("LocationStateRecent")
	var gloState globalState
	c.Find(bson.M{"_id": locationName}).One(&gloState)
	return gloState.LocationFileStatus
}

func (s mongoStore) setLocationStateRecent(locationName string, locationFileStatus map[string]string) {
	log.Println("Storing/updating recent location state")
	c := s.mongo.C("LocationStateRecent")

	var state globalState
	err := c.Find(bson.M{"_id": locationName}).One(&state)
	if err == nil {
		currentState := bson.M{"_id": locationName}
		change := bson.M{"$set": bson.M{"locationfilestatus": locationFileStatus}}
		c.Update(currentState, change)
	} else {
		state := globalState{LocationName: locationName, LocationFileStatus: locationFileStatus}
		c.Insert(state)
	}
}

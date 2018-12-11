package monitor

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"fmt"
)

type Store interface {
	setGlobalStateDaily(globalFileStatus []Location) error
	getGlobalStateDailyForThisDate(searchDate string) (Response, error)
}

type FileCheck struct {
	Location   string `json:"name,omitempty" bson:"Location,omitempty"`
	CheckDate  string `json:"checkdate,omitempty" bson:"CheckDate,omitempty"`
	FileName   string `json:"filename,omitempty" bson:"FileName,omitempty"`
	FileStatus string `json:"filestatus,omitempty" bson:"FileStatus,omitempty"`
}


type mongoStore struct {
	mongo *mgo.Database
}

type globalState struct {
	Date             string `bson:"_id,omitempty"`
	Locations []Location
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) setGlobalStateDaily(globalFileStatus []Location) error {
	log.Println("Storing daily global state")
	c := s.mongo.C("GlobalStateDaily")
	cflat := s.mongo.C("GlobalStateFlat")

	stateItem := globalState{Date: time.Now().Format("02012006"), Locations: globalFileStatus}

	var fc FileCheck
	rDate := []rune(stateItem.Date)
	fc.CheckDate = fmt.Sprintf("%s-%s-%s", string(rDate[4:8]), string(rDate[2:4]), string(rDate[0:2]))

	for _, loc := range stateItem.Locations {
		fc.Location = loc.LocationName
		for name, status := range loc.Files {
			fc.FileStatus = status
			fc.FileName = name
			log.Printf("%v\n", fc)
			err := cflat.Insert(fc)
			if err != nil {
				log.Panic(err.Error())
			}

		}
	}
	return c.Insert(stateItem)
}

func (s mongoStore) getGlobalStateDailyForThisDate(searchDate string) (Response, error) {
	log.Printf("Retreiving global state for %v", searchDate)
	c := s.mongo.C("GlobalStateDaily")
	var gs globalState
	err := c.FindId(searchDate).One(&gs)
	if err != nil {
		return Response{}, err
	}
	return Response{gs.Locations}, nil
}

package models

import (
	"github.com/globalsign/mgo/bson"
)

type Marker struct {
	X      float32       `json:"x" bson:"x"`
	Y      float32       `json:"y" bson:"y"`
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	Name   string        `json:"name" bson:"name"`
	Key    string        `json:"key" bson:"key"`
	Icon   string        `json:"icon" bson:"icon" validate:"isUrl"`
	UserId string        `json:"userId" bson:"userId"`
	Global bool          `json:"global" bson:"global"`
}

var DefaultMarker = Marker{
	X:      20,
	Y:      20,
	Id:     bson.NewObjectId(),
	Name:   "Test Marker",
	Key:    "testMarker",
	Icon:   "none",
	UserId: "testUser",
	Global: false,
}

type UserInfo struct {
	UserId     string `json:"uid"`
	PlatformId string `json:"pltid"`
	Username   string `json:"usr"`
	Level      int    `json:"level"`
	Platform   string `json:"plt"`
}

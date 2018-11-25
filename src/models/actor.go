package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_NAME_Actor  = "actor"

type Actor struct {
	Id_           bson.ObjectId `bson:"_id"`
	Name          string  `bson:"name"`
	CreateTime    time.Time     `bson:"create_time"`
	Avatar        string        `bson:"avatar"`
}

func NewActor() *Actor {
	obj := &Actor{}
	obj.CreateTime = time.Now()
	obj.Id_ = bson.NewObjectId()
	return obj
}
package models

import (
	"dao/baseSession"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const Collection_Name_Songs string = "songs"

type Songs struct {
	Id_        bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"`
	Lyricist   *mgo.DBRef    `bson:"lyricist"`
	Composer   *mgo.DBRef    `bson:"composer"`
	Singer     *mgo.DBRef    `bson:"singer"`
	Time       int32         `bson:"time"`
	CreateTime time.Time     `bson:"create_time"`
}

func NewSongs() *Songs {
	LyrRef := &mgo.DBRef{
		Database:   baseSession.DataBaseName,
		Collection: "",
	}
	newObj := &Songs{}
	newObj.CreateTime = time.Now()
	newObj.Id_ = bson.NewObjectId()
	newObj.Lyricist = LyrRef
	return newObj
}

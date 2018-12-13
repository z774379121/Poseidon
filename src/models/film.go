package models

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const Collection_Name_film = "film"

type (
	FileType  int
	Dress     int
	Topic     int
	Publishes int
)

const (
	FileType_Student   FileType = iota
	FileType_Student1  FileType = iota
	FileType_Student2  FileType = iota
	FileType_Student3  FileType = iota
	FileType_Student4  FileType = iota
	FileType_Student5  FileType = iota
	FileType_Student6  FileType = iota
	FileType_Student7  FileType = iota
	FileType_Student8  FileType = iota
	FileType_Student9  FileType = iota
	FileType_Student10 FileType = iota
)

const (
	Publishe_TokyoHot  Publishes = iota
	Publishe_TokyoHot1 Publishes = iota
	Publishe_TokyoHot2 Publishes = iota
	Publishe_TokyoHot3 Publishes = iota
	Publishe_TokyoHot4 Publishes = iota
	Publishe_TokyoHot5 Publishes = iota
)

type Film struct {
	ID_         bson.ObjectId `bson:"_id"`
	VideoCode   string
	Name        string    `bson:"name"`
	ActorRef    mgo.DBRef `bson:"actor_ref"`
	Year        int       `bson:"year"`
	Type        FileType  `bson:"type"`
	Publisher   Publishes `bson:"publisher"`
	IsPrimitive bool      `bson:"is_primitive"`
	Size        int       `bson:"size"`
	FriendShip  []string  `bson:"friend_ship"`
	Dress       Dress     `bson:"dress"`
	CreateTime  time.Time `bson:"create_time"`
	Topic       Topic     `bson:"topic"`
}

func NewFilm() *Film {
	obj := &Film{}
	obj.ID_ = bson.NewObjectId()
	obj.CreateTime = time.Now()
	obj.ActorRef.Database = baseSession.DataBaseName
	obj.ActorRef.Collection = COLLECTION_NAME_Actor
	return obj
}

package models

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const Collection_Name_Tag = "tag"

type Tag struct {
	Id_        bson.ObjectId `bson:"_id"`
	Name       string
	CreateTime time.Time
	DeleteTime time.Time
	IsDeleted  bool
	UserRef    mgo.DBRef
}

func NewTag() *Tag {
	obj := &Tag{}
	obj.CreateTime = time.Now()
	obj.Id_ = bson.NewObjectId()
	obj.UserRef.Collection = COLLECTION_NAME_User
	obj.UserRef.Database = baseSession.DataBaseName
	return obj
}

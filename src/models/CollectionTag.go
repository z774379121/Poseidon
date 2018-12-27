package models

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const Collection_Name_CollectionTag string = "collectionTag"

type CollectionTag struct {
	Id_           bson.ObjectId `bson:"_id"`
	UserRef       mgo.DBRef     `bson:"user_ref"`
	TagRef        mgo.DBRef     `bson:"tag_ref"`
	CollectionRef mgo.DBRef     `bson:"collection_ref"`
	IsDelete      bool          `bson:"is_delete"`
	CreateTime    time.Time     `bson:"create_time"`
}

func NewCollectionTag() *CollectionTag {
	obj := new(CollectionTag)
	obj.Id_ = bson.NewObjectId()
	obj.CreateTime = time.Now()
	obj.UserRef.Collection = COLLECTION_NAME_User
	obj.TagRef.Collection = Collection_Name_Tag
	obj.CollectionRef.Collection = Collection_Name_Collection
	obj.TagRef.Database = baseSession.DataBaseName
	obj.CollectionRef.Database = baseSession.DataBaseName
	obj.UserRef.Database = baseSession.DataBaseName
	return obj
}

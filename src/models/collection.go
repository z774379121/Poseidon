package models

import (
	"time"

	"github.com/z774379121/untitled1/src/dao/baseSession"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const Collection_Name_Collection = "collection"

type Shapeness int

const (
	Shapeness_NotDefine Shapeness = iota
	Shapeness_360P
	Shapeness_480P
	Shapeness_720P
	Shapeness_1080P
	Shapeness_2K
)

type Collection struct {
	Id_        bson.ObjectId `bson:"_id"`
	UserRef    mgo.DBRef     `bson:"user_ref"`
	FilmRef    mgo.DBRef     `bson:"film_ref"`
	Shapeness  Shapeness     `bson:"shapeness"`
	LocalPath  string        `bson:"local_path"`
	BTLink     string        `bson:"bt_link"`
	BTFileName string        `bson:"bt_file_name"`
	CreateTime time.Time     `bson:"create_time"`
	IsDelete   bool          `bson:"is_delete"`
	DeleteTime time.Time     `bson:"delete_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Content    string        `bson:"content"`
	Tag        *mgo.DBRef    `bson:"tag"`
	IsFavorite bool          `bson:"is_favorite"`
}

func NewColletion() *Collection {
	obj := &Collection{}
	obj.UserRef.Database = baseSession.DataBaseName
	obj.UserRef.Collection = COLLECTION_NAME_User
	obj.FilmRef.Collection = Collection_Name_film
	obj.FilmRef.Database = baseSession.DataBaseName
	obj.CreateTime = time.Now()
	obj.Id_ = bson.NewObjectId()
	return obj
}

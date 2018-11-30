package models

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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
	UserRef    mgo.DBRef
	FilmRef    mgo.DBRef
	Shapeness  Shapeness
	LocalPath  string
	BTLink     string
	BTFileName string
	CreateTime time.Time
	IsDelete   bool
	DeleteTime time.Time
	UpdateTime time.Time
	Content    string
	Tag        *mgo.DBRef
	IsFavorite bool
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

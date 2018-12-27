package dao

import (
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type DaoCollectionTag interface {
	FindByUidAndTagIds(uid bson.ObjectId, TagIds []bson.ObjectId) *models.CollectionTag
	InsertOne(m *models.CollectionTag) bool
	DeleteOneByUidAndCollectionIdAndTagId(uid, cid, tid bson.ObjectId) bool
	DeleteAllByCollectionIdAndUid(uid, cid bson.ObjectId) bool
	DeleteAllByTagIdAndUid(uid, tid bson.ObjectId) bool
}

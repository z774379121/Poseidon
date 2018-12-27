package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoCollection interface {
	SelectByUid(uid bson.ObjectId) *[]models.Collection
	SelectByTagIdsAndUid(uid bson.ObjectId, tags []bson.ObjectId) *[]models.Collection
	RemoveTagByTagIdAndUidCollectionId(tid, uid, cid bson.ObjectId) bool
	RemoveColletionTagRelation(tid, uid bson.ObjectId) bool
	DeleteOneByIdAndUid(id, uid bson.ObjectId) bool
	DeleteAllByTagIdAndUid(tid, uid bson.ObjectId) bool
	InsertOne(collection *models.Collection) bool
	AddTagByIdUidAndTagId(uid, tagId, id bson.ObjectId) bool
}

func NewDaoCollection() daoCollection {
	return daoImpl.NewDaoCollection()
}

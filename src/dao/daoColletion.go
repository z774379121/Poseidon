package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoCollection interface {
	SelectByUid(uid bson.ObjectId) *[]models.Collection
	SelectByTagIdAndUid(tid, uid bson.ObjectId) *[]models.Collection
	DeleteOneByIdAndUid(id, uid bson.ObjectId) bool
	DeleteAllByTagIdAndUid(tid, uid bson.ObjectId) bool
	InsertOne(collection *models.Collection) bool
}

func NewDaoCollection() daoCollection {
	return daoImpl.NewDaoCollection()
}

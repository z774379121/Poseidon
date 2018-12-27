package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoTag interface {
	InsertOne(tag *models.Tag) bool
	SeletTagsByUid(uid bson.ObjectId) *[]models.Tag
	DeleteOneByIdAndUid(id, uid bson.ObjectId) bool
	SeletTagByName(tagName string, uid bson.ObjectId) *models.Tag
}

func NewDaoTag() daoTag {
	return daoImpl.NewDaoTag()
}

package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoBTFileContent interface {
	SelectByName(name string) *models.BTFiles
	UpdateUserRef(urealName string, id bson.ObjectId) bool
	UpdateRealFileName(name, realName string) bool
	UpdateCollectionRef(cId bson.ObjectId, name string) bool
}

func NewDaoBTFileContent() daoBTFileContent {
	return daoImpl.NewDaoBTFileContent()
}

package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoActor interface {
	SelectById(id bson.ObjectId) *models.Actor
	TotalCount(m bson.M) int
	InsertModel(m *models.Actor) bool
	SelectByName(name string) *models.Actor
	InsertModels(m *[]models.Actor) bool
	UpdateModel(m *models.Actor) bool
	SelectByPage(page int) *[]models.Actor
	SelectLikeByName(name string) *[]models.Actor
}

var testDaoActor daoActor

func InitTestDaoActor(testInterface daoActor) {
	testDaoActor = testInterface
}

func NewDaoActor() daoActor {
	//if testDaoActor != nil {
	//	return testDaoActor
	//}
	return daoImpl.NewDaoActorImpl()
}

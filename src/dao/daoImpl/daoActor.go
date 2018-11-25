package daoImpl

import (
	"dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
	"models/modelsDefine"
	"github.com/smallnest/rpcx/log"
)

type daoActorImp struct {
	baseSession.BaseSession
}

func NewDaoActorImpl() *daoActorImp {
	dao := &daoActorImp{}
	dao.Init(models.Actor{}, models.COLLECTION_NAME_Actor)
	return dao
}

func (this *daoActorImp) SelectById(id bson.ObjectId) *models.Actor {
	result := &models.Actor{}

	err := this.FindOne(bson.M{modelsDefine.MDActor_Id_: id}, result)
	if err != nil {
		log.Error(err, id)
		return nil
	}
	return result
}

func (this *daoActorImp) TotalCount(m bson.M) int {
	count, err := this.FindCount(m)
	if err != nil {
		log.Error(err)
		return 0
	}
	return count
}

func (this *daoActorImp) SelectByName(name string) *models.Actor {
	var actor models.Actor
	err := this.FindOne(bson.M{modelsDefine.MDActor_Name: name}, &actor)
	if err != nil {
		log.Error("未找到符合该token的用户", err)
		return nil
	} else {
		return &actor
	}

}

func (this *daoActorImp) InsertModel(m *models.Actor) bool {
	return this.BaseSession.InsertModel(m)
}

func (this *daoActorImp) InsertModels(m *[]models.Actor) bool {
	return this.BaseSession.InsertManyModel(m)
}

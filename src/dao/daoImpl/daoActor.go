package daoImpl

import (
	"fmt"
	"github.com/smallnest/rpcx/log"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
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
		log.Error("未找到符合的演员", err)
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

func (this *daoActorImp) UpdateModel(m *models.Actor) bool {
	return this.BaseSession.UpdateModel(m.Id_, m)
}

func (this *daoActorImp) SelectByPage(page int) *[]models.Actor {
	return this.BaseSession.SelectByPagination(bson.M{modelsDefine.MDActor_Avatar: bson.M{baseSession.MGO_SELECT_NE: ""}}, page, 24).(*[]models.Actor)
}

func (this *daoActorImp) SelectLikeByName(name string) *[]models.Actor {
	result := this.BaseSession.SelectByPagination(bson.M{modelsDefine.MDActor_Name: bson.M{baseSession.MGO_SELECT_REGEX: fmt.Sprintf("^%s", name)}}, 0, 30)
	return result.(*[]models.Actor)
}

package daoImpl

import (
	"fmt"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
)

type daoFilmImpl struct {
	baseSession.BaseSession
}

func NewDaoFilmImpl() *daoFilmImpl {
	obj := new(daoFilmImpl)
	obj.Init(models.Film{}, models.Collection_Name_film)
	return obj
}

func (this *daoFilmImpl) FindByActorId(id bson.ObjectId) *[]models.Film {
	result := &[]models.Film{}
	if ok := this.BaseSession.SelectAll(bson.M{modelsDefine.MDFilm_ActorRef_id: id}, result); !ok {
		fmt.Println("查询失败")
		return nil
	}
	return result
}

func (this *daoFilmImpl) InsertOne(film *models.Film) bool {
	return this.BaseSession.InsertModel(film)

}

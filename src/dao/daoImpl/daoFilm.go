package daoImpl

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
)

type daoFilmImpl struct {
	baseSession.BaseSession
}

func (this *daoFilmImpl) FindByCode(code string) *models.Film {
	result := &models.Film{}
	one := this.BaseSession.FindOne(bson.M{modelsDefine.MDFilm_Name: code}, result)
	if one != nil {
		log.Info("bad find", code)
		return nil
	}
	return result
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

package daoImpl

import (
	"github.com/labstack/gommon/log"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/logger"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
)

type daoTagImpl struct {
	baseSession.BaseSession
}

func (this *daoTagImpl) SeletTagsByUid(uid bson.ObjectId) *[]models.Tag {
	result := &[]models.Tag{}
	if ok := this.BaseSession.SelectAll(bson.M{modelsDefine.MDTag_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}, result); !ok {
		log.Info("查找tag失败", uid)
	}
	return result
}

func (this *daoTagImpl) SeletTagByName(tagName string, uid bson.ObjectId) *models.Tag {
	result := &models.Tag{}
	err := this.FindOne(bson.M{modelsDefine.MDTag_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false, modelsDefine.MDTag_Name: tagName}, result)
	if err != nil {
		log.Info(err)
		return nil
	}
	return result
}

func NewDaoTag() *daoTagImpl {
	obj := new(daoTagImpl)
	obj.Init(models.Tag{}, models.Collection_Name_Tag)
	return obj

}

func (this *daoTagImpl) InsertOne(tag *models.Tag) bool {
	return this.BaseSession.InsertModel(tag)
}

func (this *daoTagImpl) DeleteOneByIdAndUid(id, uid bson.ObjectId) bool {
	return this.Update(bson.M{modelsDefine.MDTag_Id_: id, modelsDefine.MDTag_UserRef_id: uid}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDTag_IsDeleted: true}}) == nil
}

func (this *daoTagImpl) SelectCountByUid(uid bson.ObjectId) int {
	if count, err := this.FindCount(bson.M{modelsDefine.MDTag_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}); err != nil {
		logger.Sugar.Error(err)
		return 0
	} else {
		return count
	}
}

func (this *daoTagImpl) SelectByUidAndId(uid, id bson.ObjectId) *models.Tag {
	ret := &models.Tag{}
	err := this.FindOne(bson.M{modelsDefine.MDTag_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false, modelsDefine.MDTag_Id_: id}, ret)
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}
	return ret
}

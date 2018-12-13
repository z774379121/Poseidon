package daoImpl

import (
	"github.com/labstack/gommon/log"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type daoCollectionImpl struct {
	baseSession.BaseSession
}

func (this *daoCollectionImpl) DeleteOneByIdAndUid(id, uid bson.ObjectId) bool {
	return this.BaseSession.Update(bson.M{modelsDefine.MDCollection_Id_: id, modelsDefine.MDCollection_UserRef_id: uid}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDCollection_IsDelete: true, modelsDefine.MDCollection_DeleteTime: time.Now()}}) == nil
}

func (this *daoCollectionImpl) DeleteAllByTagIdAndUid(tid, uid bson.ObjectId) bool {
	return this.BaseSession.UpdateAll(bson.M{modelsDefine.MDCollection_Tag_id: tid, modelsDefine.MDCollection_UserRef_id: uid}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDCollection_IsDelete: true, modelsDefine.MDCollection_DeleteTime: time.Now()}}) == nil
}

func (this *daoCollectionImpl) SelectByUid(uid bson.ObjectId) *[]models.Collection {

	results := &[]models.Collection{}
	if !this.BaseSession.SelectAll(bson.M{modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}, results) {
		log.Info("查看标签下收藏失败", uid)
	}
	return results
}

func (this *daoCollectionImpl) SelectByTagIdAndUid(tid, uid bson.ObjectId) *[]models.Collection {
	results := &[]models.Collection{}
	if !this.BaseSession.SelectAll(bson.M{modelsDefine.MDCollection_Tag_id: tid, modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}, results) {
		log.Info("查看标签下收藏失败", tid, uid)
	}
	return results
}

func (this *daoCollectionImpl) InsertOne(collection *models.Collection) bool {
	return this.BaseSession.InsertModel(collection)
}

func NewDaoCollection() *daoCollectionImpl {
	obj := new(daoCollectionImpl)
	obj.Init(models.Collection{}, models.Collection_Name_Collection)
	return obj
}

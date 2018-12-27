package daoImpl

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/logger"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type daoCollectionImpl struct {
	baseSession.BaseSession
}

// 根据集合id删除某个集合.
func (this *daoCollectionImpl) DeleteOneByIdAndUid(id, uid bson.ObjectId) bool {
	return this.BaseSession.Update(bson.M{modelsDefine.MDCollection_Id_: id, modelsDefine.MDCollection_UserRef_id: uid}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDCollection_IsDelete: true, modelsDefine.MDCollection_DeleteTime: time.Now()}}) == nil
}

// 删除所有带有x标签的集合(删除标签后调用该方法).
func (this *daoCollectionImpl) DeleteAllByTagIdAndUid(tid, uid bson.ObjectId) bool {
	return this.BaseSession.UpdateAll(bson.M{modelsDefine.MDCollection_Tags: tid, modelsDefine.MDCollection_UserRef_id: uid}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDCollection_IsDelete: true, modelsDefine.MDCollection_DeleteTime: time.Now()}}) == nil
}

// 查询某用户的所有集合
func (this *daoCollectionImpl) SelectByUid(uid bson.ObjectId) *[]models.Collection {
	results := &[]models.Collection{}
	if !this.BaseSession.SelectAll(bson.M{modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}, results) {
		logger.Sugar.Info("查看标签下收藏失败", uid)
	}
	return results
}

// 查询含有所有所需标签的集合
func (this *daoCollectionImpl) SelectByTagIdsAndUid(uid bson.ObjectId, tags []bson.ObjectId) *[]models.Collection {
	results := &[]models.Collection{}
	if !this.BaseSession.SelectAll(bson.M{modelsDefine.MDCollection_Tags: bson.M{
		baseSession.MGO_SELECT_All: tags,
	}, modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDTag_IsDeleted: false}, results) {
		logger.Sugar.Info("查看标签下收藏失败", tags, uid)
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

// 为某集合增加一个标签
func (this *daoCollectionImpl) AddTagByIdUidAndTagId(uid, tagId, id bson.ObjectId) bool {
	err := this.Update(bson.M{modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDCollection_Id_: id}, bson.M{baseSession.MGO_UPDATE_push: bson.M{modelsDefine.MDCollection_Tags: tagId}})
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err == nil
}

// 移除某集合的x标签.
func (this *daoCollectionImpl) RemoveTagByTagIdAndUidCollectionId(tid, uid, cid bson.ObjectId) bool {
	err := this.Update(bson.M{modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDCollection_Id_: cid}, bson.M{baseSession.MGO_UPDATE_pull: bson.M{modelsDefine.MDCollection_Tags: tid}})
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err == nil
}

// 去掉所有包含x标签集合上的x标签(不删除集合).
func (this *daoCollectionImpl) RemoveColletionTagRelation(tid, uid bson.ObjectId) bool {
	err := this.Update(bson.M{modelsDefine.MDCollection_UserRef_id: uid, modelsDefine.MDCollection_Tags: tid}, bson.M{baseSession.MGO_UPDATE_pull: bson.M{modelsDefine.MDCollection_Tags: tid}})
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err == nil
}

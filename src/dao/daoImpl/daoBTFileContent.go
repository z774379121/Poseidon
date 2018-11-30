package daoImpl

import (
	"fmt"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
)

type daoBTFileContentImpl struct {
	baseSession.BaseSession
}

func NewDaoBTFileContent() *daoBTFileContentImpl {
	obj := &daoBTFileContentImpl{}
	obj.Init(models.BTFiles{}, models.COLLECTION_NAME_BTFileContent)
	return obj
}

func (this *daoBTFileContentImpl) UpdateUserRef(name string, uid bson.ObjectId) bool {
	err := this.Update(bson.M{modelsDefine.MDBTFiles_FileName: name}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDBTFiles_UserRef_id: uid}})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (this *daoBTFileContentImpl) UpdateRealFileName(name, realName string) bool {
	err := this.Update(bson.M{modelsDefine.MDBTFiles_FileName: name}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDBTFiles_RealFileName: realName}})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (this *daoBTFileContentImpl) UpdateCollectionRef(cId bson.ObjectId, name string) bool {
	err := this.Update(bson.M{modelsDefine.MDBTFiles_FileName: name}, bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDBTFiles_CollectionRef_id: cId}})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func (this *daoBTFileContentImpl) SelectByName(name string) *models.BTFiles {
	result := &models.BTFiles{}
	err := this.FindOne(bson.M{modelsDefine.MDBTFiles_FileName: name}, result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

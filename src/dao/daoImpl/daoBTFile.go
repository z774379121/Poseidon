package daoImpl

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoBTFileImpl struct {
	baseSession.BaseSession
}

func NewdaoBTFileImpl() *daoBTFileImpl {
	obj := &daoBTFileImpl{}
	obj.Init(models.BTFile{}, models.COLLECTION_NAME_BTFile)
	return obj
}
func (this *daoBTFileImpl) UploadBTFile(uid bson.ObjectId, filedata *[]byte) (bool, string) {
	fileName := bson.NewObjectId().Hex()
	this.CreateFile(fileName, *filedata)
	return true, fileName
}
func (this *daoBTFileImpl) DownloadBTFile(filename string) *[]byte {
	return this.ReadFile(filename)
}

package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"gopkg.in/mgo.v2/bson"
)

type daoBTFile interface {
	UploadBTFile(uid bson.ObjectId, filedata *[]byte) (bool, string)
	DownloadBTFile(filename string) *[]byte
}

var testDaoBTFile daoBTFile

func NewDaoBTFile() daoBTFile {
	if testDaoBTFile != nil {
		return testDaoBTFile
	}
	return daoImpl.NewDaoBTFile()
}

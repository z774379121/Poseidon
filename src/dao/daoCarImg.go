package dao

import (
	"dao/daoImpl"
	"gopkg.in/mgo.v2/bson"
)

type daoCarImg interface {
	UploadImg(uid bson.ObjectId, filedata *[]byte) (bool, string)
	DownloadImg(filename string) *[]byte
}

var testDaoCarImg daoCarImg

func NewDaoCarImg() daoCarImg {
	if testDaoCarImg != nil {
		return testDaoCarImg
	}
	return daoImpl.NewDaoCarImgImpl()
}

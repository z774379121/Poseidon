package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_NAME_CardCouponImg = "cardCoupon"

type CardCouponImg struct {
	Id_        bson.ObjectId `bson:"_id"`
	ChunkSize  int           `bson:"chunkSize"`
	UploadDate time.Time     `bson:"uploadDate"`
	Length     int           `bson:"length"`
	Md5        string        `bson:"md5"`
	FileName   string        `bson:"filename"`

	SongRef mgo.DBRef `bson:"user_ref"`
}

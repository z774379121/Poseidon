package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_NAME_Actor = "actor"
const SHOW_NUM_PERPAGE_Actor = 24

type Cup int

const (
	ACup Cup = iota
	BCup
	CCup
	DCup
	ECup
	FCup
	GCup
	HCup
)

func (this Cup) toString() string {
	switch this {
	case ACup:
		return "ACup"
	case BCup:
		return "BCup"
	case CCup:
		return "CCup"
	case DCup:
		return "DCup"
	case ECup:
		return "ECup"
	case FCup:
		return "FCup"
	case GCup:
		return "GCup"
	case HCup:
		return "HCup"
	default:
		return ""
	}
}

type Actor struct {
	Id_              bson.ObjectId `bson:"_id"`
	Name             string        `bson:"name"`
	CreateTime       time.Time     `bson:"create_time"`
	Avatar           string        `bson:"avatar"`
	BirthDay         time.Time     `bson:"birth_day"`
	BirthPalce       string        `bson:"birth_palce"`
	Height           int           `bson:"height"`
	Habit            string        `bson:"habit"`
	Cup              Cup           `bson:"cup"`
	Bust             int           `bson:"bust"`
	WaistLine        int           `bson:"waist_line"`
	HipCircumference int           `bson:"hip_circumference"`
}

func NewActor() *Actor {
	obj := &Actor{}
	obj.CreateTime = time.Now()
	obj.Id_ = bson.NewObjectId()
	return obj
}

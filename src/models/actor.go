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

var Cmap = map[string]Cup{
	"A": ACup,
	"B": BCup,
	"C": CCup,
	"D": DCup,
	"E": ECup,
	"F": FCup,
	"G": GCup,
	"H": HCup,
}

func (this Cup) toString() string {
	switch this {
	case ACup:
		return "A"
	case BCup:
		return "B"
	case CCup:
		return "C"
	case DCup:
		return "D"
	case ECup:
		return "E"
	case FCup:
		return "F"
	case GCup:
		return "G"
	case HCup:
		return "H"
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

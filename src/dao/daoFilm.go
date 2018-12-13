package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoFilm interface {
	FindByActorId(id bson.ObjectId) *[]models.Film
	InsertOne(film *models.Film) bool
}

func NewDaoFilm() daoFilm {
	return daoImpl.NewDaoFilmImpl()
}

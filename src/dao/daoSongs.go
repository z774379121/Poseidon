package dao

import (
	"dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
)

type daoSongs interface {
	InsertOne(songs *models.Songs) bool
	SelectByName(name string) *models.Songs
}

var testDaoSongs daoSongs

func NewDaoSongs() daoSongs {
	if testDaoSongs != nil {
		return testDaoSongs
	}
	return daoImpl.NewDaoSongsImpl()
}

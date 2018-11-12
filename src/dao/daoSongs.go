package dao

import (
	"dao/daoImpl"
	"models"
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

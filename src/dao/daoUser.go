package dao

import (
	"dao/daoImpl"
	"gopkg.in/mgo.v2/bson"
	"models"
)

type daoUser interface {
	SelectByUId(id bson.ObjectId) *models.User
	TotalCount(m bson.M) int
	UpdateSessionAndToken(id bson.ObjectId, session, token string) bool
	SelectByEmail(email string) *models.User
	InsertModel(m *models.User) bool
	SelectByAppToken(appToken string) *models.User
	UpdateUserPassword(id bson.ObjectId, newPassword string) bool
}

var testDaoUser daoUser

func InitTestDaoUser(testInterface daoUser) {
	testDaoUser = testInterface
}

func NewDaoUser() daoUser {
	if testDaoUser != nil {
		return testDaoUser
	}
	return daoImpl.NewDaoUserImpl()
}

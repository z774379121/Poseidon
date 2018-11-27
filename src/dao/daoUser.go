package dao

import (
	"github.com/z774379121/untitled1/src/dao/daoImpl"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoUser interface {
	SelectByUId(id bson.ObjectId) *models.User
	TotalCount(m bson.M) int
	UpdateSessionAndToken(id bson.ObjectId, session, token string) bool
	SelectByEmail(email string) *models.User
	SelectByEmailAll(email string) *models.User
	InsertModel(m *models.User) bool
	SelectByAppToken(appToken string) *models.User
	UpdateUserPassword(id bson.ObjectId, newPassword string) bool
	UpdateUserEmailCheck(id bson.ObjectId) bool
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

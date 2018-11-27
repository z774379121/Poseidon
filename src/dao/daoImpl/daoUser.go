package daoImpl

import (
	"github.com/smallnest/rpcx/log"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/models/modelsDefine"
	"gopkg.in/mgo.v2/bson"
)

type daoUserImp struct {
	baseSession.BaseSession
}

func NewDaoUserImpl() *daoUserImp {
	dao := &daoUserImp{}
	dao.Init(models.User{}, models.COLLECTION_NAME_User)
	return dao
}

func (this *daoUserImp) SelectByUId(id bson.ObjectId) *models.User {
	result := &models.User{}

	err := this.FindOne(bson.M{modelsDefine.MDUser_Id_: id}, result)
	if err != nil {
		log.Error(err, id)
		return nil
	}
	return result
}

func (this *daoUserImp) SelectByEmail(email string) *models.User {
	result := &models.User{}

	err := this.FindOne(bson.M{modelsDefine.MDUser_Email: email, modelsDefine.MDUser_IsConfirmed: true}, result)
	if err != nil {
		log.Error(err, email)
		return nil
	}
	return result
}

func (this *daoUserImp) SelectByEmailAll(email string) *models.User {
	result := &models.User{}

	err := this.FindOne(bson.M{modelsDefine.MDUser_Email: email}, result)
	if err != nil {
		log.Error(err, email)
		return nil
	}
	return result
}

func (this *daoUserImp) TotalCount(m bson.M) int {
	count, err := this.FindCount(m)
	if err != nil {
		log.Error(err)
		return 0
	}
	return count
}

func (this *daoUserImp) SelectByAppToken(appToken string) *models.User {
	var user models.User
	err := this.FindOne(bson.M{modelsDefine.MDUser_Session: appToken[:32], modelsDefine.MDUser_Token: appToken[32:]}, &user)
	if err != nil {
		log.Error("未找到符合该token的用户", err)
		return nil
	} else {
		return &user
	}

}

func (this *daoUserImp) UpdateSessionAndToken(id bson.ObjectId, session, token string) bool {
	err := this.Update(bson.M{modelsDefine.MDUser_Id_: id},
		bson.M{
			baseSession.MGO_UPDATE_SET: bson.M{
				modelsDefine.MDUser_Session: session,
				modelsDefine.MDUser_Token:   token,
			},
		},
	)
	if err != nil {
		log.Error(err, id)
		return false
	}
	return true
}

func (this *daoUserImp) InsertModel(m *models.User) bool {
	return this.BaseSession.InsertModel(m)
}

func (this *daoUserImp) UpdateUserPassword(id bson.ObjectId, newPassword string) bool {
	if err := this.Update(bson.M{modelsDefine.MDUser_Id_: id},
		bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDUser_Password: newPassword}},
	); err != nil {
		log.Error(err, id)
		return false
	} else {
		return true
	}
}

func (this *daoUserImp) UpdateUserEmailCheck(id bson.ObjectId) bool {
	if err := this.Update(bson.M{modelsDefine.MDUser_Id_: id},
		bson.M{baseSession.MGO_UPDATE_SET: bson.M{modelsDefine.MDUser_IsConfirmed: true}},
	); err != nil {
		log.Error(err, id)
		return false
	} else {
		return true
	}
}

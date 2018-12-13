package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/dao"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func NewFilm(context echo.Context) error {
	f := new(FilmStruct)
	if err := context.Bind(f); err != nil {
		log.Println("表单无效", err)
	}
	token := context.Get("token")

	//daoUser := dao.NewDaoUser()
	//user := daoUser.SelectByAppToken(token.(string))
	////daoClt := dao.NewColletion()
	//daoFile := dao.NewDaoBTFileContent()
	//daoFile.SelectByName(f.BTFileName)
	//
	//newObj := models.NewColletion()
	////newObj.FilmRef.Id_ = f.FilmName
	//newObj.Content = f.Content
	//
	//newObj.BTLink = f.BTLink
	//newObj.LocalPath = f.LocalPath
	//newObj.Shapeness = models.Shapeness(f.Shapeness)
	//newObj.UserRef.Id = user.Id_
	fmt.Println(token)
	fmt.Println(f)
	return context.String(http.StatusOK, "ok")
}

func GetColletions(ctx echo.Context) error {
	token := ctx.Get("token")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(token.(string))
	if user == nil {
		return ctx.String(http.StatusBadRequest, "bad uid")
	}

	daoCollection := dao.NewDaoCollection()
	collections := daoCollection.SelectByUid(user.Id_)
	if collections != nil {
		return ctx.String(http.StatusBadRequest, "无效查询")
	}

	return ctx.JSON(http.StatusOK, collections)
}

func GetCollectionUnderTag(ctx echo.Context) error {
	token := ctx.Get("token")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(token.(string))
	if user == nil {
		return ctx.String(http.StatusBadRequest, "bad uid")
	}
	tag := ctx.QueryParam("tag")
	if !bson.IsObjectIdHex(tag) {
		return ctx.String(http.StatusBadRequest, "bad tagId")
	}
	daoCollection := dao.NewDaoCollection()
	collections := daoCollection.SelectByTagIdAndUid(bson.ObjectIdHex(tag), user.Id_)
	return ctx.JSON(http.StatusOK, collections)
}

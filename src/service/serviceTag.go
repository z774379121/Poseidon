package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/controller"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type FSTags struct {
	Name       string    `json:"name"`
	Id         string    `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

func GetTags(ctx echo.Context) error {
	uid := ctx.Get("token").(string)
	user := controller.GetUser(uid)
	daoTag := dao.NewDaoTag()
	tags := daoTag.SeletTagsByUid(user.Id_)
	if tags == nil {
		return ctx.String(http.StatusBadRequest, "查询tags失败")
	}
	results := make([]*FSTags, len(*tags))
	for idx, tag := range *tags {
		obj := &FSTags{}
		obj.Name = tag.Name
		obj.Id = tag.Id_.Hex()
		obj.CreateTime = tag.CreateTime
		results[idx] = obj
	}
	return ctx.JSON(http.StatusOK, results)
}

func GetTagDetail(ctx echo.Context) error {
	tagName := ctx.Param("tname")
	uid := ctx.Get("token").(string)
	user := controller.GetUser(uid)
	daoTag := dao.NewDaoTag()
	tag := daoTag.SeletTagByName(tagName, user.Id_)
	return ctx.JSON(http.StatusOK, tag)
}

type Ftag struct {
	Name string `json:"name"`
}

func NewTag(ctx echo.Context) error {
	newFtag := &Ftag{}
	err := ctx.Bind(newFtag)
	if err != nil {
		return err
	}
	uid := ctx.Get("token").(string)
	user := controller.GetUser(uid)
	daoTag := dao.NewDaoTag()
	fmt.Println(newFtag.Name)
	count := daoTag.SelectCountByUid(user.Id_)
	if count >= 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "最多同时拥有10个标签")
	}
	tag := daoTag.SeletTagByName(newFtag.Name, user.Id_)
	fmt.Println(tag)
	if tag != nil {
		return ctx.String(http.StatusBadRequest, "tag已经存在")
	}
	newTag := models.NewTag()
	newTag.Name = newFtag.Name
	newTag.UserRef.Id = user.Id_
	ok := daoTag.InsertOne(newTag)
	if !ok {
		return ctx.String(http.StatusBadRequest, "插入失败")
	}
	return ctx.Redirect(http.StatusFound, "/user/tags")
}

func DeleteOneTag(ctx echo.Context) error {
	tag := ctx.Param("tname")
	fmt.Println(tag)
	if !bson.IsObjectIdHex(tag) {
		return echo.ErrValidatorNotRegistered
	}
	uid := ctx.Get("token").(string)
	user := controller.GetUser(uid)
	daoTag := dao.NewDaoTag()
	target := daoTag.SelectByUidAndId(user.Id_, bson.ObjectIdHex(tag))
	if target == nil {
		return echo.ErrNotFound
	}
	if done := daoTag.DeleteOneByIdAndUid(bson.ObjectIdHex(tag), user.Id_); !done {
		return echo.NewHTTPError(http.StatusBadRequest, "删除失败")
	}
	return ctx.String(http.StatusOK, "删除成功")

}

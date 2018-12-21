package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func Boot(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	daoActor := dao.NewDaoActor()
	pageInt, err := strconv.Atoi(page)
	count := daoActor.TotalCount(nil)
	if (count+models.SHOW_NUM_PERPAGE_Actor-1)/models.SHOW_NUM_PERPAGE_Actor < pageInt {
		return ctx.String(http.StatusForbidden, "page over")
	}
	if err != nil {
		pageInt = 0
	}
	actors := daoActor.SelectByPage(pageInt)
	return ctx.Render(http.StatusOK, "bs3.html", map[string]interface{}{
		"users": actors,
		"change": func(row int) bool {
			return row%6 == 0 && row != 0
		},
	})
}

func Bootvue(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	daoActor := dao.NewDaoActor()
	pageInt, err := strconv.Atoi(page)
	count := daoActor.TotalCount(nil)
	if (count+models.SHOW_NUM_PERPAGE_Actor-1)/models.SHOW_NUM_PERPAGE_Actor < pageInt {
		return ctx.String(http.StatusForbidden, "page over")
	}
	if err != nil {
		pageInt = 0
	}
	actors := daoActor.SelectByPage(pageInt)
	return ctx.JSON(http.StatusOK, actors)
}
func BootLike(ctx echo.Context) error {
	name := ctx.FormValue("name")
	daoActor := dao.NewDaoActor()
	actors := daoActor.SelectLikeByName(name)
	return ctx.Render(http.StatusOK, "bs3.html", map[string]interface{}{
		"users": actors,
		"change": func(row int) bool {
			return row%6 == 0 && row != 0
		},
	})
}
func GetActor(ctx echo.Context) error {
	actorId := ctx.Param("id")
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.Response().Header().Set("content-type", "application/json")
	fmt.Println(actorId)
	if !bson.IsObjectIdHex(actorId) {
		return ctx.String(http.StatusForbidden, "bad Id")
	}
	daoActor := dao.NewDaoActor()
	actor := daoActor.SelectById(bson.ObjectIdHex(actorId))
	if actor == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}
	return ctx.JSON(http.StatusOK, actor)
}

func GetActorRegex(ctx echo.Context) error {
	queryStr := ctx.QueryParam("like")
	daoActor := dao.NewDaoActor()
	actors := daoActor.SelectLikeByName(queryStr)
	return ctx.JSON(http.StatusOK, actors)
}

func GetActorByHeadChar(context echo.Context) error {
	name := context.QueryParam("keyword")
	daoActor := dao.NewDaoActor()
	actors := daoActor.SelectLikeByName(name)
	data := len(*actors)
	if data > 10 {
		data = 10
	}
	z := make([]string, 0, data)
	actor := *actors
	for i := 0; i < data; i++ {
		actor := actor[i]
		z = append(z, actor.Name)
	}
	return context.JSON(http.StatusOK, map[string]interface{}{
		"s":    z,
		"data": true,
	})
}

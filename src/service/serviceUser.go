package service

import (
	"controller"
	"dao"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"models"
	"net/http"
	"setting"
	"time"
)

func SignUp(context echo.Context) error {
	userName := context.FormValue("userName")
	pwd := context.FormValue("passWord")
	email := context.FormValue("email")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByEmail(email)
	if user != nil {
		return context.String(http.StatusForbidden, "邮箱已经被使用")
	}
	user = models.NewUser()
	user.Password = pwd
	user.Email = email
	user.GenSalt()
	user.EncodePasswd()
	insertModel := daoUser.InsertModel(user)
	if !insertModel {
		return context.String(http.StatusBadRequest, "插入新用户到数据库失败")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(setting.JWTSignKey))
	if err != nil {
		return err
	}
	return context.String(http.StatusOK, fmt.Sprintf("http://localhost:1323/auth/?token=%s", t))
}

func Login(ctx echo.Context) error {
	pwd := ctx.FormValue("passWord")
	email := ctx.FormValue("email")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByEmail(email)
	if user == nil {
		return ctx.String(http.StatusForbidden, "邮箱未注册或未进行验证")
	}
	if user.ValidatePassword(pwd) {
		user.GenUserToken()
		ctl := controller.NewBaseController(ctx)
		ctl.SetToken("user", user.GetAppToken())
		ctl.SetCookies()
		if !daoUser.UpdateSessionAndToken(user.Id_, user.Session, user.Token) {
			return ctx.String(http.StatusInternalServerError, "err")
		}
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	} else {
		return ctx.String(http.StatusUnauthorized, "密码错误")
	}
}

func Hello(ctx echo.Context) error {
	ctx.Response().After(func() {
		fmt.Println("after return")
	})
	ctx.Response().Before(func() {
		fmt.Println("before return")
	})
	return ctx.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}

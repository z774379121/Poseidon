package service

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/smallnest/rpcx/log"
	"github.com/z774379121/untitled1/src/controller"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"github.com/z774379121/untitled1/src/setting"
)

type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func SignUp(context echo.Context) error {
	userName := context.FormValue("username")
	pwd := context.FormValue("password")
	pwd1 := context.FormValue("password1")
	email := context.FormValue("emailsignup")
	fmt.Println(pwd, pwd1)
	if pwd1 != pwd {
		return context.HTML(http.StatusBadRequest, "<script>alert('密码不一致');</script>")
	}
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByEmail(email)
	if user != nil {
		return context.String(http.StatusForbidden, "邮箱已经被使用")
	}
	fmt.Println(userName, email)
	user = daoUser.SelectByEmailAll(email)
	if user != nil {
		return context.String(http.StatusUnauthorized, "账号已经注册,请前往邮箱确认")
	}
	user = models.NewUser()
	user.Username = userName
	user.Password = pwd
	user.Email = email
	user.GenSalt()
	user.EncodePasswd()
	insertModel := daoUser.InsertModel(user)
	if !insertModel {
		return context.String(http.StatusBadRequest, "插入新用户到数据库失败")
	}

	// Set claims
	claims := &JwtCustomClaims{
		userName,
		email,
		false,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(setting.JWTSignKey))
	if err != nil {
		return err
	}
	return context.String(http.StatusOK, fmt.Sprintf("http://localhost:1323/auth/?token=%s", t))
}

func Login(ctx echo.Context) error {
	pwd := ctx.FormValue("password")
	email := ctx.FormValue("email")
	daoUser := dao.NewDaoUser()
	userWithoutCheck := daoUser.SelectByEmailAll(email)
	if userWithoutCheck == nil {
		return ctx.String(http.StatusForbidden, "邮箱不存在")
	}
	if !userWithoutCheck.IsConfirmed {
		return ctx.String(http.StatusUnauthorized, "邮箱未确认")
	}
	if userWithoutCheck.ValidatePassword(pwd) {
		userWithoutCheck.GenUserToken()
		ctl := controller.NewBaseController(ctx)
		ctl.SetToken("user", userWithoutCheck.GetAppToken())
		ctl.SetCookies()
		if !daoUser.UpdateSessionAndToken(userWithoutCheck.Id_, userWithoutCheck.Session, userWithoutCheck.Token) {
			return ctx.String(http.StatusInternalServerError, "err")
		}
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	} else {
		return ctx.String(http.StatusUnauthorized, "密码错误")
	}
}

func Logout(ctx echo.Context) error {
	ctl := controller.NewBaseController(ctx)
	ctl.ClearCookies()
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(ctl.GetUserToken())
	if user == nil {
		return ctx.String(http.StatusForbidden, "非法token")
	}
	ok := daoUser.UpdateSessionAndToken(user.Id_, "", "")
	if !ok {
		return ctx.String(http.StatusInternalServerError, "清空token失败")
	}
	return ctx.String(http.StatusOK, "logout")
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

func DP(ctx echo.Context) error {
	if baseSession.DumpData() == nil {
		return ctx.String(http.StatusOK, "pong")
	}
	return ctx.String(http.StatusInternalServerError, "down")
}

func LoginH(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index.html", nil)
}

func SignUpH(context echo.Context) error {
	return context.Redirect(http.StatusMovedPermanently, "/login#toregister")
}

func C(context echo.Context) error {
	if err := baseSession.DumpData(); err != nil {
		log.Fatal(err)
	}
	daoActor := dao.NewDaoActor()

	oo := models.NewActor()
	oo.Name = "test"
	b := daoActor.InsertModel(oo)
	if !b {
		log.Fatal("插入失败")
	}
	fi, err := os.Open("C:/Users/JJ/Desktop/code/ag.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	lines := [103]*[]models.Actor{}
	for i := 0; i < len(lines); i++ {
		lines[i] = &[]models.Actor{}
	}
	count := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		avatar_name := strings.Split(string(a), "#")
		avatar, name := avatar_name[0], avatar_name[1]
		obj := models.NewActor()
		obj.Name = name
		obj.Avatar = avatar
		index := count / 100
		*lines[index] = append(*lines[index], *obj)
		count += 1
	}
	fmt.Println(len(lines))
	for i := 0; i < len(lines); i++ {
		insertModels := daoActor.InsertModels(lines[i])
		if !insertModels {
			log.Fatal("插入失败")
		}
		fmt.Println("插入成功", i)
	}
	return nil
}

// 邮箱链接确认,通过jwt
func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	name := claims.Name
	email := claims.Email
	daoUser := dao.NewDaoUser()
	muser := daoUser.SelectByEmailAll(email)
	if !daoUser.UpdateUserEmailCheck(muser.Id_) {
		return c.String(http.StatusOK, "认证失败")
	}

	return c.String(http.StatusOK, "Welcome "+name+"!")
}

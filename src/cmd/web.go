package cmd

import (
	"bytes"
	"controller"
	"dao"
	"dao/baseSession"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"
	"io"
	"github.com/z774379121/untitled1/src/models"
	"net/http"
	"os"
	"service"
	"setting"
	"time"
	"html/template"
)

var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `Gogs web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "port, p", Value: "3000", Usage: "Temporary port number to prevent conflict"},
		cli.StringFlag{Name: "config, c", Value: "/src/config/cfg.ini", Usage: "Custom configuration file path"},
	},
}

func runWeb(context *cli.Context) error {
	CfgFile := context.String("config")
	if context.IsSet("config") {
		fmt.Println("custom file:", CfgFile)
	}
	setting.CfgFileName = CfgFile
	setting.GlobalInit()
	baseSession.DBInit()
	t := &service.TemplateRenderer{
		Templates: template.Must(template.ParseGlob("src/view/*.html")),
	}
	e := echo.New()
	e.Static("/", "src/view")
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", service.Hello)
	e.GET("/dp", service.DP)
	e.POST("/signUp", service.SignUp)
	e.GET("/signUp", service.SignUpH)
	e.POST("/login", service.Login)
	e.GET("/login", service.LoginH)
	//e.GET("/c", service.C)
	admin := e.Group("/admin")
	{
		admin.Use(ServiceController)
		admin.GET("/", func(context echo.Context) error {
			return context.String(http.StatusOK, "welcome, admin")
		})
		admin.GET("/logout", service.Logout)
	}
	l := e.Group("/auth")
	{
		l.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey:  []byte(setting.JWTSignKey),
			TokenLookup: "query:token",
		}))
		l.GET("/", restricted)
	}
	g := e.Group("/v1")
	{
		g.GET("/", func(context echo.Context) error {
			cookie, err := context.Cookie("User_token")
			var name string
			if err != nil {
				fmt.Println(err)
				cookie := new(http.Cookie)
				cookie.Name = "User_token"
				cookie.Value = "jon"
				cookie.Expires = time.Now().Add(24 * time.Hour)
				context.SetCookie(cookie)
				cookie = new(http.Cookie)
				cookie.Name = "_tokenType"
				cookie.Value = "Manner"
				cookie.Expires = time.Now().Add(24 * time.Hour)
				context.SetCookie(cookie)
				fmt.Println("设置成功")
			} else {
				name = cookie.Value
			}

			return context.String(http.StatusOK, fmt.Sprintf("hello from web %s", name))
		})
		g.POST("/songs", func(context echo.Context) error {
			name := context.FormValue("name")
			singer := context.FormValue("singer")
			fmt.Println(singer)
			daoSongs := dao.NewDaoSongs()
			newSong := models.NewSongs()
			newSong.Name = name
			if ok := daoSongs.InsertOne(newSong); !ok {
				context.String(http.StatusInternalServerError, "插入失败")
			}
			return context.String(http.StatusOK, "插入成功")
		})
		g.POST("/upload", func(context echo.Context) error {
			avatar, err := context.FormFile("avatar")
			if err != nil {
				return err
			}

			// Source
			fmt.Println(avatar.Filename, avatar.Size)
			src, err := avatar.Open()
			if err != nil {
				return err
			}

			defer src.Close()

			des := make([]byte, 1000000)
			n, err2 := src.Read(des)
			if err2 != nil {
				fmt.Println(err2)
				return context.String(http.StatusBadRequest, "读取失败")
			}
			fmt.Println(n)
			// Destination
			dao := dao.NewDaoCarImg()
			ok, name := dao.UploadImg(bson.NewObjectId(), &des)
			if !ok {
				return context.String(http.StatusBadRequest, "插入到数据库失败")
			}
			return context.String(http.StatusOK, name)

		})
		g.GET("/download/:filename", func(context echo.Context) error {
			filename := context.Param("filename")
			dao := dao.NewDaoCarImg()
			img := dao.DownloadImg(filename)
			dst, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, bytes.NewReader(*img)); err != nil {
				return err
			}
			return context.File(filename)

		})
	}
	e.Logger.Fatal(e.Start(setting.Port))
	return nil

}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	email := claims["email"].(string)
	daoUser := dao.NewDaoUser()
	muser := daoUser.SelectByEmailAll(email)
	if !daoUser.UpdateUserEmailCheck(muser.Id_) {
		return c.String(http.StatusOK, "认证失败")
	}

	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func ServiceController(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		ctl := new(controller.BaseController)
		ctl.C = context
		token := ctl.GetUserToken()
		if token == "" {
			return ctl.C.String(http.StatusUnauthorized, "请先登录")
		}
		daoUser := dao.NewDaoUser()
		user := daoUser.SelectByAppToken(token)
		if user == nil {
			return ctl.C.String(http.StatusUnauthorized, "非法token")
		}
		return next(context)
	}
}

package cmd

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"github.com/z774379121/untitled1/src/controller"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/service"
	"github.com/z774379121/untitled1/src/setting"
)

var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `Power web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "port, p", Value: "3000", Usage: "Temporary port number to prevent conflict"},
		cli.StringFlag{Name: "config, c", Value: "/src/config/cfg.ini", Usage: "Custom configuration file path"},
	},
}

func runWeb(context *cli.Context) error {
	// Global init
	CfgFile := context.String("config")
	if context.IsSet("config") {
		fmt.Println("custom file:", CfgFile)
	}
	setting.CfgFileName = CfgFile
	setting.GlobalInit()
	baseSession.DBInit()

	// set static like echo templeRender and 404 page etc.
	t := &service.TemplateRenderer{
		Templates: template.Must(template.ParseGlob("src/view/*.html")),
	}
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.File("src/view/404.html")
	}
	//{
	//	daoFilm := dao.NewDaoFilm()
	//	daoActor := dao.NewDaoActor()
	//	actor := daoActor.SelectByName("小川桃果")
	//	obj := models.NewFilm()
	//	obj.Name = "xxxx"
	//	obj.Size = 1212
	//	obj.ActorRef.Id = actor.Id_
	//	ok := daoFilm.InsertOne(obj)
	//	if !ok {
	//		fmt.Println("插入失败")
	//	}
	//}
	e := echo.New()
	e.Static("/", "src/view")
	e.File("/favicon.ico", "/images/play.ico")
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// begin router
	e.GET("/", service.Hello)
	e.GET("/dp", service.DP)
	e.POST("/signUp", service.SignUp)
	e.GET("/signUp", service.SignUpH)
	e.POST("/login", service.Login)
	e.GET("/login", service.LoginH)
	e.GET("/boot", service.Boot)
	e.GET("/actor/:id", service.GetActor)
	e.GET("/actor", service.GetActorRegex)
	e.GET("/ss", func(context echo.Context) error {
		return context.Render(http.StatusOK, "sear.html", nil)
	})
	e.GET("/search", func(context echo.Context) error {
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
	})
	e.GET("/actors/:name", service.FindActorFilm)
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
			Claims:      &service.JwtCustomClaims{},
			SigningKey:  []byte(setting.JWTSignKey),
			TokenLookup: "query:token",
		}))
		l.GET("/", service.Restricted)
	}

	g := e.Group("/v1")
	{
		g.GET("/upload", func(context echo.Context) error {
			return context.Render(http.StatusOK, "upload.html", nil)
		}, ServiceController)
		//		g.GET("/", func(context echo.Context) error {
		//			cookie, err := context.Cookie("User_token")
		//			var name string
		//			if err != nil {
		//				fmt.Println(err)
		//				cookie := new(http.Cookie)
		//				cookie.Name = "User_token"
		//				cookie.Value = "jon"
		//				cookie.Expires = time.Now().Add(24 * time.Hour)
		//				context.SetCookie(cookie)
		//				cookie = new(http.Cookie)
		//				cookie.Name = "_tokenType"
		//				cookie.Value = "Manner"
		//				cookie.Expires = time.Now().Add(24 * time.Hour)
		//				context.SetCookie(cookie)
		//				fmt.Println("设置成功")
		//			} else {
		//				name = cookie.Value
		//			}
		//
		//			return context.String(http.StatusOK, fmt.Sprintf("hello from web %s", name))
		//		})

		g.POST("/upload", service.UpLoad)
		g.GET("/download/:filename", service.Download)
		g.POST("/fs", service.NewFilm, ServiceController)
	}
	u := e.Group("/user", ServiceController)
	{
		u.GET("/tags", service.GetTags)
		u.GET("/tag/:tname", service.GetTagDetail)
		u.POST("/tag", service.NewTag)

		u.GET("/Collections", service.GetColletions)
		u.GET("/Collection", service.GetCollectionUnderTag)
	}
	e.Logger.Fatal(e.Start(setting.Port))
	return nil

}

// 登录中间件
func ServiceController(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		ctl := new(controller.BaseController)
		ctl.C = context
		token := ctl.GetUserToken()
		if token == "" {
			return ctl.C.Redirect(http.StatusFound, "/login")
		}

		daoUser := dao.NewDaoUser()
		user := daoUser.SelectByAppToken(token)
		if user == nil {
			return context.String(http.StatusUnauthorized, "非法token")
		}

		context.Set("token", token)
		return next(context)
	}
}

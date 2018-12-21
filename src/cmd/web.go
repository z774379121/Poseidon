package cmd

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"github.com/z774379121/untitled1/src/controller"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/service"
	"github.com/z774379121/untitled1/src/setting"
	"html/template"
	"net/http"
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
	e := echo.New()
	e.Static("/", "src/view")
	e.File("/favicon.ico", "/images/play.ico")
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// begin router
	e.GET("/", func(context echo.Context) error {
		return context.File("src/view/main.html")
	})
	e.GET("/dp", service.DP)
	e.POST("/signUp", service.SignUp)
	e.GET("/signUp", service.SignUpH)
	e.POST("/login", service.Login)
	e.GET("/login", service.LoginH)
	// actorList and single info
	e.GET("/boot", service.Boot)
	e.POST("/bootlike", service.BootLike)
	// actorDetail.
	e.POST("/actor/:id", service.GetActor)
	e.GET("/actor/:id", func(context echo.Context) error {
		return context.File("src/view/actordetail.html")
	})
	e.GET("/actor", service.GetActorRegex)
	e.GET("/film", service.FindHasCode)
	e.GET("/actorvue", service.Bootvue)
	//e.GET("/tvue", func(context echo.Context) error {
	//	return context.File("src/view/bs3vue.html")
	//})
	api := e.Group("/api/vi")
	{
		api.GET("/actors/:name", service.FindActorFilm)
		api.GET("/afilms/:aid", service.GetActorFilms)
		api.GET("/actor/:id", service.GetActor)
		api.GET("/film", service.FindHasCode)
		api.GET("/search", service.GetActorByHeadChar)
	}
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
	e.GET("/films/:name", service.FindActorFilm)
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
		g.GET("/fs", func(context echo.Context) error {
			return context.Render(http.StatusOK, "films.html", nil)
		})
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

//var cc *opencc.OpenCC
//
//func Conver(s string) string {
//	nText, err := cc.ConvertText(s)
//	if err != nil {
//		fmt.Println(err)
//		return s
//	}
//	return nText
//}
//func ins() {
//	fi, err := os.Open("/home/jj/Desktop/ag.txt")
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		return
//	}
//	defer fi.Close()
//	fmt.Println(fi.Name())
//	cc, err = opencc.NewOpenCC("t2s")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	br := bufio.NewReader(fi)
//	count := 0
//	actors := make([]models.Actor, 10167)
//	for {
//		a, _, c := br.ReadLine()
//		if c == io.EOF {
//			break
//		}
//		actor := models.NewActor()
//		allInfo := string(a)
//		slice1 := strings.Split(allInfo, "|")
//		actor.Avatar = slice1[0]
//		other := slice1[1:]
//		slice2 := strings.Split(strings.Join(other, ""), "##")
//		name := slice2[0]
//		other2 := slice2[1:]
//		if len(other2) == 1 {
//
//		}
//		actor.Name = Conver(name)
//		detail := strings.Split(strings.Join(other2, ""), "#")
//		for _, value := range detail {
//			//fmt.Println(value)
//			if strings.HasPrefix(value, "生日: ") {
//				//fmt.Println("f")
//
//				birthDay := value[strings.Index(value, " ")+1:]
//				btime, err := time.Parse("2006-01-02", birthDay)
//				if err != nil {
//					fmt.Println(err)
//				}
//				actor.BirthDay = btime
//			} else if strings.HasPrefix(value, "身高: ") {
//				fmt.Println("身高", value[strings.Index(value, " ")+1:])
//				actor.Height, _ = strconv.Atoi(value[strings.Index(value, " ")+1 : len(value)-2])
//			} else if strings.HasPrefix(value, "出生地: ") {
//				fmt.Println("出生地", value[strings.Index(value, " ")+1:])
//				actor.BirthPalce = value[strings.Index(value, " ")+1:]
//			} else if strings.HasPrefix(value, "腰围: ") {
//				fmt.Println("腰围", value[strings.Index(value, " ")+1:])
//				actor.WaistLine, _ = strconv.Atoi(value[strings.Index(value, " ")+1 : len(value)-2])
//			} else if strings.HasPrefix(value, "爱好: ") {
//				fmt.Println("爱好", value[strings.Index(value, " ")+1:])
//				actor.Habit = value[strings.Index(value, " ")+1:]
//			} else if strings.HasPrefix(value, "胸围: ") {
//				fmt.Println("胸围", value[strings.Index(value, " ")+1:])
//				actor.Bust, _ = strconv.Atoi(value[strings.Index(value, " ")+1 : len(value)-2])
//			} else if strings.HasPrefix(value, "臀围: ") {
//				fmt.Println("臀围", value[strings.Index(value, " ")+1:])
//				actor.HipCircumference, _ = strconv.Atoi(value[strings.Index(value, " ")+1 : len(value)-2])
//			} else if strings.HasPrefix(value, "罩杯: ") {
//				fmt.Println("罩杯", value[strings.Index(value, " ")+1:])
//				actor.Cup = models.Cmap[value[strings.Index(value, " ")+1:]]
//			}
//		}
//		actors[count] = *actor
//		count += 1
//	}
//	fmt.Println(count)
//	daoactor := dao.NewDaoActor()
//	insertModels := daoactor.InsertModels(&actors)
//	fmt.Println(insertModels)
//}

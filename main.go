package main

import (
	"github.com/z774379121/untitled1/src/cmd"
	"github.com/urfave/cli"
	"os"
	"setting"
)

const APP_VER = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = setting.ServiceName
	app.Usage = "A painless self-hosted Git service"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.BackUp,
		cmd.Restore,
	}
	app.Run(os.Args)
}

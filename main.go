package main

import (
	"./src/cmd"
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
	}
	app.Run(os.Args)
}

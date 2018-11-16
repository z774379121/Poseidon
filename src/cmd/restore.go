package cmd

import (
	"dao/baseSession"
	"fmt"
	"github.com/mongodb/mongo-tools/common/db"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongorestore"
	"github.com/urfave/cli"
	"setting"
	"strings"
)

var Restore = cli.Command{
	Name:  "restore",
	Usage: "Restore files and database from backup",
	Description: `Restore imports all related files and database from a backup archive.
The backup version must lower or equal to current Gogs version. You can also import
backup from other database engines, which is useful for database migrating.

If corresponding files or database tables are not presented in the archive, they will
be skipped and remain unchanged.`,
	Action: runRestore,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "/src/config/cfg.ini", Usage: "Custom configuration file path"},
	},
}

func runRestore(ctx *cli.Context) {
	CfgFile := ctx.String("config")
	if ctx.IsSet("config") {
		fmt.Println("custom file:", CfgFile)
	}

	setting.CfgFileName = CfgFile
	setting.GlobalInit()
	baseSession.DBInit()
	inputOptions := &mongorestore.InputOptions{}
	outputOptions := &mongorestore.OutputOptions{
		NumParallelCollections: 1,
		NumInsertionWorkers:    1,
	}
	nsOptions := &mongorestore.NSOptions{}
	hostPort := strings.Split(setting.DBConfig.Host, ":")
	connection := &options.Connection{
		Host: hostPort[0],
		Port: "27018",
	}
	auth := options.Auth{
		Username:  setting.DBConfig.UserName,
		Password:  setting.DBConfig.Password,
		Source:    setting.DBConfig.DatabaseName,
		Mechanism: "SCRAM-SHA-1",
	}
	toolOptions := &options.ToolOptions{
		SSL: &options.SSL{
			UseSSL: false,
		},
		Connection: connection,
		Auth:       &auth,
		Verbosity:  &options.Verbosity{},
		URI:        &options.URI{},
	}
	toolOptions.Namespace = &options.Namespace{DB: setting.DBConfig.DatabaseName}
	sessionProvider, err := db.NewSessionProvider(*toolOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	restore := mongorestore.MongoRestore{
		ToolOptions:     toolOptions,
		OutputOptions:   outputOptions,
		InputOptions:    inputOptions,
		NSOptions:       nsOptions,
		SessionProvider: sessionProvider,
	}
	restore.TargetDirectory = "/home/jj/go/src/github.com/mongodb/mongo-tools/mongodump/dump"
	err = restore.Restore()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("opk")
}

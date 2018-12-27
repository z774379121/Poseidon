package cmd

import (
	"fmt"
	"github.com/Unknwon/cae/zip"
	"github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
	"github.com/urfave/cli"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/logger"
	"github.com/z774379121/untitled1/src/setting"
	"io/ioutil"
	"os"
	path2 "path"
	"strings"
	"time"
)

var BackUp = cli.Command{
	Name:  "backup",
	Usage: "Backup database and file",
	Description: `Backup dumps and compresses all related files and database into zip file,
which can be used for migrating current server to another server. The output format is meant to be
portable among all supported database engines.`,
	Action: backup,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "archive-name", Usage: "Name of backup archive", Value: fmt.Sprintf("gogs-backup-%s.zip", time.Now().Format("20060102150405"))},
		cli.StringFlag{Name: "target", Usage: "Target directory path to save backup archive", Value: "dump"},
		cli.StringFlag{Name: "config, c", Value: "/src/config/cfg.ini", Usage: "Custom configuration file path"},
	},
}

func backup(ctx *cli.Context) {
	CfgFile := ctx.String("config")
	//target := ctx.String("target")
	archive := ctx.String("archive-name")
	fmt.Println(archive)
	if ctx.IsSet("config") {
		fmt.Println("custom file:", CfgFile)
	}

	setting.CfgFileName = CfgFile
	setting.GlobalInit()
	baseSession.DBInit()
	auth := options.Auth{
		Username:  setting.DBConfig.UserName,
		Password:  setting.DBConfig.Password,
		Source:    setting.DBConfig.DatabaseName,
		Mechanism: "SCRAM-SHA-1",
	}
	// 默认备份当前使用的数据库源
	hostPort := strings.Split(setting.DBConfig.Host, ":")
	connection := &options.Connection{
		Host: hostPort[0],
		Port: hostPort[1],
	}
	toolOptions := &options.ToolOptions{
		SSL: &options.SSL{
			UseSSL: false,
		},
		Connection: connection,
		Auth:       &auth,
		Verbosity:  &options.Verbosity{},
	}

	toolOptions.Namespace = &options.Namespace{DB: setting.DBConfig.DatabaseName}

	outputOptions := &mongodump.OutputOptions{
		NumParallelCollections: 1,
	}
	inputOptions := &mongodump.InputOptions{}

	log.SetVerbosity(toolOptions.Verbosity)

	md := &mongodump.MongoDump{
		ToolOptions:   toolOptions,
		InputOptions:  inputOptions,
		OutputOptions: outputOptions,
	}
	tmpDir, err := ioutil.TempDir("/tmp", "tmpDir-")
	fmt.Println("建立临时目录:", tmpDir)
	db := path2.Join(tmpDir, "db")
	md.OutputOptions.Out = db
	err = md.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = md.Dump()
	if err != nil {
		fmt.Println(err)
		return
	}

	archivez, err := zip.Create(archive)
	if err != nil {
		fmt.Println(err)
		return
	}
	//err = archivez.AddFile("sd.jpg", "/home/jj/Desktop/untitled1/5bebddcab91c2037c121188f.jpg")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	err = archivez.AddDir("gog/database", db)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = archivez.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	os.RemoveAll(tmpDir)
	logger.Sugar.Info("ok", archive)

	return
}

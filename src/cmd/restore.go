package cmd

import (
	"fmt"
	"github.com/Unknwon/cae/zip"
	"github.com/mongodb/mongo-tools/common/db"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongorestore"
	"github.com/urfave/cli"
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/logger"
	"github.com/z774379121/untitled1/src/setting"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
		cli.StringFlag{Name: "file, f", Usage: "restore file"},
	},
}

func runRestore(ctx *cli.Context) {
	CfgFile := ctx.String("config")
	if ctx.IsSet("config") {
		fmt.Println("custom file:", CfgFile)
	}
	// cfg init
	setting.CfgFileName = CfgFile
	setting.GlobalInit()
	baseSession.DBInit()

	// unzip file
	var zipFile string
	if ctx.IsSet("file") {
		zipFile = ctx.String("file")
	} else {
		zipFile = searchBackupFiles()
		if zipFile == "" {
			logger.Sugar.Error("没有找到备份文件,请检查当前目录或指定对应文件")
			return
		}
	}
	logger.Sugar.Info("将使用备份源:", zipFile)
	tmpDir := os.TempDir()
	zip.ExtractTo(zipFile, tmpDir)
	archiveDir := path.Join(tmpDir, "gog")
	defer os.RemoveAll(archiveDir)

	// restore mongodb
	dbFile := path.Join(archiveDir, "database")
	inputOptions := &mongorestore.InputOptions{}
	outputOptions := &mongorestore.OutputOptions{
		NumParallelCollections: 1,
		NumInsertionWorkers:    1,
	}
	nsOptions := &mongorestore.NSOptions{}
	connection := &options.Connection{
		Host: "127.0.0.1",
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
		logger.Sugar.Error(err)
		return
	}
	restore := mongorestore.MongoRestore{
		ToolOptions:     toolOptions,
		OutputOptions:   outputOptions,
		InputOptions:    inputOptions,
		NSOptions:       nsOptions,
		SessionProvider: sessionProvider,
	}
	restore.TargetDirectory = dbFile
	err = restore.Restore()
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Infof("成功导入备份数据到 数据库: %s", setting.DBConfig.DatabaseName)
}

// 若未指定备份文件, 将使用当前目录下最新的.zip备份源文件
func searchBackupFiles() (ret string) {
	zips, err := filepath.Glob("*.zip")
	if err != nil {
		logger.Sugar.Info(err)
		return
	}
	if len(zips) == 0 {
		return
	}
	recently := 0
	index := 0
	for idx, z := range zips {
		if len(z) == len("gogs-backup-20060102150405.zip") {
			loopt, _ := strconv.Atoi(z[12:26])
			if loopt > recently {
				recently = loopt
				index = idx
			}
		}
	}
	ret = zips[index]
	return ret
}

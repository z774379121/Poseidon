package setting

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

var (
	CfgFileName             string
	Cfg                     *ini.File
	Port                    string
	DevModel                bool
	MongoDBMaxConnectionNum int
	ServiceName             string
	DBConfig                struct {
		DatabaseName string
		UserName     string
		Password     string
		Host         string
	}
)

func GlobalInit() {
	pwd, _ := os.Getwd()
	Cfg, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true}, pwd+CfgFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	Cfg.NameMapper = ini.AllCapsUnderscore
	sec := Cfg.Section("security")
	Port = sec.Key("SERVICE_PORT").MustString(":1323")
	DevModel = sec.Key("IS_DEV_MODE").MustBool(false)
	MongoDBMaxConnectionNum = sec.Key("MONGODB_MAX_CONNECTION_NUM").MustInt(300)
	ServiceName = sec.Key("SERVICE_NAME").MustString("unKnow")
	if DevModel {
		sec = Cfg.Section("DB")
	} else {
		sec = Cfg.Section("DBTianYi")
	}
	err = sec.MapTo(&DBConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}

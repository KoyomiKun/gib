package util

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	LogLevel  string
	LogFormat string

	JwtKey string

	AccessKey string
	SecretKey string
	Bucket    string
	Server    string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误:", err)
	}
	loadServer(file)
	loadDb(file)
	loadLogger(file)
	loadMiddleware(file)
	loadQiNiu(file)
}

func loadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":9000")
}
func loadMiddleware(file *ini.File) {
	JwtKey = file.Section("middleware").Key("JwtKey").MustString("")
}
func loadQiNiu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").MustString("")
	Bucket = file.Section("qiniu").Key("Bucket").MustString("")
	SecretKey = file.Section("qiniu").Key("SecretKey").MustString("")
	Server = file.Section("qiniu").Key("Server").MustString("")

}

func loadDb(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassword = file.Section("database").Key("DbPassword").MustString("123456")
	DbName = file.Section("database").Key("DbName").MustString("ginblog")
}

func loadLogger(file *ini.File) {
	LogLevel = file.Section("logger").Key("LogLevel").MustString("debug")
	LogFormat = file.Section("logger").Key("LogFormat").MustString("2006-01-02 15:04:05")
}

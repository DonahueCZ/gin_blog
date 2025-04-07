package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

// 引入config的参数作为全局变量
var (
	AppMode  string
	HttpPort string
	JwtKey   string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	AccessKey  string
	SecretKey  string
	Bucket     string
	QiniuSever string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取出错", err)
	}
	LoadServer(file)
	LoadDatabase(file)
	LoadQiniu(file)
}

func LoadServer(file *ini.File) {
	//ini包的用法，Section指定了读取的模块，Key指定了读取具体的值，如果值存在，就会读取键值，如果值不存在读取MustString的默认值
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("87wqho45wdiq56")

}

func LoadDatabase(file *ini.File) {
	Db = file.Section("server").Key("Db").MustString("mysql")
	DbHost = file.Section("server").Key("DbHost").MustString("localhost")
	DbPort = file.Section("server").Key("DbPort").MustString("3306")
	DbUser = file.Section("server").Key("DbUser").MustString("ginblog")
	DbPassWord = file.Section("server").Key("DbPassWord").MustString("admin123")
	DbName = file.Section("server").Key("DbName").MustString("ginblog")
}

func LoadQiniu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SecretKey = file.Section("qiniu").Key("SecretKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuSever = file.Section("qiniu").Key("QiniuSever").String()
}

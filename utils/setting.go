package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

const (
	QiniuStorageType = "qiniu"
	MinioStorageType = "minio"
)

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

	StorageType string

	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string
	StorageSever     string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
	//LoadQiniu(file)
	LoadObjectStorageType(file)
	LoadObjectStorage(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("ginblog")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("admin123")
	DbName = file.Section("database").Key("DbName").MustString("ginblog")
}

func LoadObjectStorageType(file *ini.File) {
	StorageType = file.Section("storage").Key("Type").String()
}

func LoadObjectStorage(file *ini.File) {
	StorageAccessKey = file.Section("storage").Key("AccessKey").String()
	StorageSecretKey = file.Section("storage").Key("SecretKey").String()
	StorageBucket = file.Section("storage").Key("Bucket").String()
	StorageSever = file.Section("storage").Key("Sever").String()
}

//func LoadQiniu(file *ini.File) {
//	AccessKey = file.Section("qiniu").Key("AccessKey").String()
//	SecretKey = file.Section("qiniu").Key("SecretKey").String()
//	Bucket = file.Section("qiniu").Key("Bucket").String()
//	QiniuSever = file.Section("qiniu").Key("QiniuSever").String()
//}

//func LoadMinio() {
//
//}
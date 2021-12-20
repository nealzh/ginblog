package utils

import (
	"crypto/rsa"
	"fmt"
	utils "ginblog/utils/rsa"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
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

	StorageAccessKey         string
	StorageSecretKey         string
	StorageBucket            string
	StorageSever             string
	StorageExpirationSeconds int

	RsaPublicKey  *rsa.PublicKey
	RsaPrivateKey *rsa.PrivateKey
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}

	LoadRsaKey()

	LoadServer(file)
	LoadData(file)
	LoadObjectStorageType(file)
	LoadObjectStorage(file)
}

func LoadFileContent(filePath string) []byte {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return fileContent
}

func LoadRsaKey() {
	RsaPublicKey = utils.BytesToPublicKey(LoadFileContent("config/public.key"))
	RsaPrivateKey = utils.BytesToPrivateKey(LoadFileContent("config/private.key"))
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
	StorageExpirationSeconds, _ = file.Section("storage").Key("ExpirationSeconds").Int()
}

package model

import (
	"context"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"ginblog/utils/string_util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
	"log"
	"time"

	//"github.com/minio/minio-go/v7"
	//"github.com/minio/minio-go/v7/pkg/credentials"
	//"log"
	"mime/multipart"
)

var StorageType = utils.StorageType

var AccessKey = utils.StorageAccessKey
var SecretKey = utils.StorageSecretKey
var Bucket = utils.StorageBucket
var ServerUrl = utils.StorageSever

type Object struct {
	User User `gorm:"foreignkey:Uid"`
	gorm.Model
	Uid         int       `gorm:"type:int;not null"`
	Name        string    `gorm:"type:varchar(128)"`
	Suffix      string    `gorm:"type:varchar(128)" json:"suffix"`
	ContentType string    `gorm:"type:varchar(128)" json:"content_type"`
	URL         string    `gorm:"type:text" json:"object_url"`
	Expiration  time.Time `gorm:"type:datetime(3)" json:"expiration"`
}

func GetDownloadUrl(oid int) (string, string, int) {
	return "", "", errmsg.ERROR
}

func UpLoadFile(userid string, file multipart.File, fileSize int64, fileSuffix string, fileContentType string) (string, string, int) {

	switch StorageType {

	case utils.QiniuStorageType:
		return UpLoadQiniuFile(file, fileSize)

	case utils.MinioStorageType:
		return UpLoadMinioFile(userid, file, fileSize, fileSuffix, fileContentType)

	default:
		return "", "", errmsg.ERROR
	}
}

func UpLoadMinioFile(userid string, file multipart.File, fileSize int64, fileSuffix string, fileContentType string) (string, string, int) {

	ctx := context.Background()

	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(ServerUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return "", "", errmsg.ERROR
	}

	//objName := userid + "/" + string_util.GenUUIDStr(32, string_util.GlobalGenBaseStr)
	objName := string_util.GenUUIDStr(32, string_util.GlobalGenBaseStr)

	if len(fileSuffix) > 0 {
		objName = objName + "." + fileSuffix
	}

	if len(fileContentType) == 0 {
		fileContentType = "application/octet-stream"
	}

	uploadInfo, err := minioClient.PutObject(ctx, Bucket, objName, file, fileSize, minio.PutObjectOptions{ContentType: fileContentType})

	if err != nil {
		log.Fatalln(err)
		return "", "", errmsg.ERROR
	}

	log.Println(uploadInfo)

	return objName, fileContentType, errmsg.SUCCSE
}

func UpLoadQiniuFile(file multipart.File, fileSize int64) (string, string, int) {
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", "", errmsg.ERROR
	}
	url := ServerUrl + ret.Key
	return url, "", errmsg.SUCCSE

}

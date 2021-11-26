package v1

import (
	"context"
	"ginblog/model"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"ginblog/utils/string_util"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"mime/multipart"
	"reflect"
	"strconv"

	//"log"
	"net/http"
	"strings"
)

var StorageType = utils.StorageType

var AccessKey = utils.StorageAccessKey
var SecretKey = utils.StorageSecretKey
var Bucket = utils.StorageBucket
var ServerUrl = utils.StorageSever

func GetDownloadUrl(c *gin.Context) {

	oid, _ := strconv.Atoi(c.Param("oid"))

	url, contentType, code := model.GetDownloadUrl(oid)

	c.JSON(http.StatusOK, gin.H{
		"status":      code,
		"message":     errmsg.GetErrMsg(code),
		"url":         url,
		"contentType": contentType,
	})
}

// UpLoad 上传图片接口
func UpLoad(c *gin.Context) {

	//userid, err := c.Request.Cookie("userid")

	var objName string
	var contentType string
	var code int

	//if err != nil {
	//
	//	url , code = "", errmsg.ERROR
	//
	//} else {
	//
	//	file, fileHeader, _ := c.Request.FormFile("file")
	//
	//	fileSize := fileHeader.Size
	//
	//	url, code = model.UpLoadObject(userid.Value, file, fileSize)
	//}

	file, fileHeader, _ := c.Request.FormFile("file")
	fileSuffixSlice := strings.Split(fileHeader.Filename, ".")
	fileSuffixSliceLen := len(fileSuffixSlice)

	//log.Println("Content-Type: ", fileHeader.Header.Get("Content-Type"))
	fileContentType := fileHeader.Header.Get("Content-Type")

	var fileSuffix string

	if fileSuffixSliceLen >= 2 {
		fileSuffix = fileSuffixSlice[fileSuffixSliceLen-1]
	} else {
		fileSuffix = ""
	}

	fileSize := fileHeader.Size

	session := sessions.Default(c)

	log.Println("userid: ", session.Get("userid"), reflect.TypeOf(session.Get("userid")))

	userid := session.Get("userid").(uint)

	log.Println("userid: ", userid)

	objName, contentType, code = UpLoadObject(userid, file, fileSize, fileSuffix, fileContentType)

	objectData := model.Object{}
	objectData.Uid = userid
	objectData.Name = objName
	objectData.Suffix = fileSuffix
	objectData.ContentType = contentType

	oid, _ := model.CreateObject(&objectData)

	c.JSON(http.StatusOK, gin.H{
		"status":      code,
		"message":     errmsg.GetErrMsg(code),
		"url":         objName,
		"contentType": contentType,
		"oid":         oid,
	})
}

func UpLoadObject(userid uint, file multipart.File, fileSize int64, fileSuffix string, fileContentType string) (string, string, int) {

	switch StorageType {

	case utils.QiniuStorageType:
		return UpLoadQiniuObject(file, fileSize)

	case utils.MinioStorageType:
		return UpLoadMinioObject(userid, file, fileSize, fileSuffix, fileContentType)

	default:
		return "", "", errmsg.ERROR
	}
}

func UpLoadMinioObject(userid uint, file multipart.File, fileSize int64, fileSuffix string, fileContentType string) (string, string, int) {

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

	objName := strconv.FormatUint(uint64(userid), 10) + "/" + string_util.GenUUIDStr(32, string_util.GlobalGenBaseStr)
	//objName := string_util.GenUUIDStr(32, string_util.GlobalGenBaseStr)

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

func UpLoadQiniuObject(file multipart.File, fileSize int64) (string, string, int) {
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

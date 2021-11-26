package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"strconv"

	//"log"
	"net/http"
	"strings"
)

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

	var url string
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
	//	url, code = model.UpLoadFile(userid.Value, file, fileSize)
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

	url, contentType, code = model.UpLoadFile("", file, fileSize, fileSuffix, fileContentType)

	c.JSON(http.StatusOK, gin.H{
		"status":      code,
		"message":     errmsg.GetErrMsg(code),
		"url":         url,
		"contentType": contentType,
	})
}

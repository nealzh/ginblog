package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
	"time"
)

type Object struct {
	User User `gorm:"foreignkey:Uid"`
	gorm.Model
	Uid         uint      `gorm:"type:int;not null"`
	Name        string    `gorm:"type:varchar(128)"`
	Suffix      string    `gorm:"type:varchar(128)" json:"suffix"`
	ContentType string    `gorm:"type:varchar(128)" json:"content_type"`
	URL         string    `gorm:"type:text" json:"object_url"`
	Expiration  time.Time `gorm:"type:datetime(3)" json:"expiration"`
}

func CreateObject(data *Object) (uint, int) {
	err := db.Create(&data).Error
	if err != nil {
		return 0, errmsg.ERROR // 500
	}
	return data.ID, errmsg.SUCCSE
}

func GetDownloadUrl(oid int) (string, string, int) {
	return "", "", errmsg.ERROR
}

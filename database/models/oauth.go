package models

import (
	"github.com/jinzhu/gorm"

	"github.com/moos3/gin-rest-api/lib/common"
)


type OauthToken struct {
	gorm.Model
	Token string
	AppID uint
	User   User `gorm:"foreignkey:UserID"`
	UserID uint
}


func (o *OauthToken) Read(m common.JSON) {
	o.ID = uint(m["id"].(float64))
	o.Token = m["token"].(string)
	o.AppID = m["app_id"].(uint)
	o.UserID = m["user_id"].(uint)
}

func (o *OauthToken) Serialize() common.JSON{
	return common.JSON{
		"id":     o.ID,
		"token":  o.Token,
		"app_id": o.AppID,
	}
}


func (o *OauthToken) Validate(m common.JSON) {

}

func (o *OauthToken) Refresh(m common.JSON){

}

func (o *OauthToken) NewOauthAccessToken(m common.JSON){

}
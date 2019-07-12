package models

import (
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/lib/common"
)

type User struct {
	gorm.Model
	Username     string
	DisplayName  string
	PasswordHash string
}

func (u *User) Serialize() common.JSON {
	return common.JSON{
		"id":           u.ID,
		"username":     u.Username,
		"display_name": u.DisplayName,
	}
}

func (u *User) Read(m common.JSON) {
	u.ID = unit(m["id"].(float64))
	u.Username = m["username"].(string)
	u.DisplayName = m["display_name"].(string)
}

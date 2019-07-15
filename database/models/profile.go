package models

import (
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/lib/common"
)

// UserProfile user profile
type UserProfile struct {
	gorm.Model
	Bio            string `sql:"type:text"`
	User           User   `gorm:"foreignkey:UserID"`
	UserID         uint
	TwitterHandle  string
	FirstName      string
	LastName       string
	AvatarImageURL string
}

// Serialize serialize profile data
func (u UserProfile) Serialize() common.JSON {
	return common.JSON{
		"uuid":         u.User,
		"bio":          u.Bio,
		"twitter":      u.TwitterHandle,
		"first_name":   u.FirstName,
		"last_name":    u.LastName,
		"avatar_image": u.AvatarImageURL,
	}
}

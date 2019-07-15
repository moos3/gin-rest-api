package models

import (
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/lib/common"
)

// User data model
type User struct {
	gorm.Model
	Username     string
	Email        string
	DisplayName  string
	PasswordHash string
	Region       string
	Verified     bool
}

// ResetPasswordToken -Password Reset Tokens
type ResetPasswordToken struct {
	gorm.Model
	Token         string
	Expiration    int64
	User          User `gorm:"foreignkey:UserID"`
	UserID        uint
	Claimed       bool
	RequestedByIP string
	UsedByIP      string
}

// JwtToken - This is so we can disable tokens
type JwtToken struct {
	gorm.Model
	TokenSha   string
	UserID     User `gorm:"foreignkey:UserID"`
	Deactivate bool
}

// Serialize serializes jwt token record data
func (j *JwtToken) Serialize() common.JSON {
	return common.JSON{
		"id":          j.ID,
		"token_sig":   j.TokenSha,
		"deactivated": j.Deactivate,
	}
}

// Serialize serializes reset token data
func (r *ResetPasswordToken) Serialize() common.JSON {
	return common.JSON{
		"token":           r.Token,
		"expiration":      r.Expiration,
		"claimed":         r.Claimed,
		"requested_by_ip": r.RequestedByIP,
		"used_by_ip":      r.UsedByIP,
	}
}

// Serialize serializes user data
func (u *User) Serialize() common.JSON {
	return common.JSON{
		"id":           u.ID,
		"username":     u.Username,
		"display_name": u.DisplayName,
		"region":       u.Region,
		"email":        u.Email,
	}
}

func (u *User) Read(m common.JSON) {
	u.ID = uint(m["id"].(float64))
	u.Username = m["username"].(string)
	u.DisplayName = m["display_name"].(string)
	u.Region = m["region"].(string)
	u.Email = m["email"].(string)
}

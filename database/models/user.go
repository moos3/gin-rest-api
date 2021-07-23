package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/google/uuid"
	"github.com/moos3/gin-rest-api/lib/common"
)

// User data model
type User struct {
	ID           string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username     string
	Email        string
	DisplayName  string
	PasswordHash string
	Region       string
	Verified     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// ResetPasswordToken -Password Reset Tokens
type ResetPasswordToken struct {
	ID            string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Token         string
	Expiration    int64
	User          User   `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID        string `gorm:"type:uuid REFERENCES users(id)"`
	Claimed       bool
	RequestedByIP string
	UsedByIP      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

// JwtToken - This is so we can disable tokens
type JwtToken struct {
	ID         string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	TokenSha   string
	User       User      `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID     string `gorm:"type:uuid REFERENCES users(id)"`
	Deactivate bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
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

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("ID", uuid)
}

func (u *User) Read(m common.JSON) {
	u.ID = m["id"].(string)
	u.Username = m["username"].(string)
	u.DisplayName = m["display_name"].(string)
	u.Region = m["region"].(string)
	u.Email = m["email"].(string)
}

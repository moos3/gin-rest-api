package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/moos3/gin-rest-api/lib/common"
)

// UserProfile user profile
type UserProfile struct {
	ID             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Bio            string    `sql:"type:text"`
	User           User      `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID         uuid.UUID `gorm:"type:uuid REFERENCES users(id)"`
	TwitterHandle  string
	FirstName      string
	LastName       string
	AvatarImageURL string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
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

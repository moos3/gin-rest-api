package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/lib/common"
)

// Token data model
type Token struct {
	gorm.Model
	Token  string
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
}

// Serialize serializes user data
func (t *Token) Serialize() common.JSON {
	return common.JSON{
		"id":    t.ID,
		"token": t.Token,
	}
}

func (t *Token) Read(m common.JSON) {
	t.ID = uint(m["id"].(float64))
	t.Token = m["token"].(string)
}

// GenerateToken - Creates a token for a user.
func (t *Token) GenerateToken() {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Failed to make token")
	}
	t.Token = id.String()
}

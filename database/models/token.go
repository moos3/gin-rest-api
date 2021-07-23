package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/moos3/gin-rest-api/lib/common"
)

// Token data model
type Token struct {
	ID        string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Token     string
	User      User   `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID    string `gorm:"type:uuid REFERENCES users(id)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Serialize serializes user data
func (t *Token) Serialize() common.JSON {
	return common.JSON{
		"id":    t.ID,
		"token": t.Token,
	}
}

func (t *Token) Read(m common.JSON) {
	t.ID = m["id"].(string)
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

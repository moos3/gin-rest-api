package models

import (
	"time"
)

// Exchange data model
type Exchange struct {
	ID          string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string
	ApiEndpoint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// UserExchange - data model
type UserExchange struct {
	ID         string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	User       User      `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID     string `gorm:"type:uuid REFERENCES users(id)"`
	Exchange   Exchange  `gorm:"ForeginKey:ExchangeID;AssociationForeignKey:ID"`
	ExchangeID string `gorm:"type:uuid REFERENCES exchanges(id)"`
	ApiKey     string
	ApiToken   string
	UsedByIP   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

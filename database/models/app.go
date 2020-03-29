package models

import "github.com/moos3/gin-rest-api/lib/common"

type App struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
	AppURL       string `json:"app_url"`
	ID           uint   `json:"id,omitempty"`
	Token        string `json:"token,omitempty"`
	AppID        uint   `json:"app_id,omitempty"`
	UserID       uint   `json:"user_id,omitempty"`
}

func (a *App) Read(m common.JSON) {
	a.ID = uint(m["id"].(float64))
	a.Token = m["token"].(string)
	a.AppID = m["app_id"].(uint)
	a.UserID = m["user_id"].(uint)
}

func (a *App) Serialize() common.JSON {
	return common.JSON{
		"id":     a.ID,
		"token":  a.Token,
		"app_id": a.AppID,
	}
}
